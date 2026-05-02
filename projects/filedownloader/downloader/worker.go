package downloader

import (
	"context"
	"errors"
	"sync"
	"time"
)

type WorkerPool interface {
	SubmitJob(ctx context.Context, job Job) error
	Close()
	Wait()
}

type workerPool struct {
	client *Client
	wg     sync.WaitGroup
	jobCh  chan Job

	retry *RetryHandler

	ctx    context.Context
	cancel context.CancelFunc

	mu      sync.Mutex
	closing bool
}

func NewWorkerPool(parent context.Context, client *Client, workers int, buffer int) WorkerPool {
	ctx, cancel := context.WithCancel(parent)

	wp := &workerPool{
		client: client,
		jobCh:  make(chan Job, buffer),
		ctx:    ctx,
		cancel: cancel,
	}

	// init retry handler (IMPORTANT)
	wp.retry = NewRetryHandler(wp)

	// start workers
	for i := 0; i < workers; i++ {
		wp.wg.Add(1)
		go func() {
			defer wp.wg.Done()
			wp.worker()
		}()
	}

	return wp
}

func (wp *workerPool) worker() {
	for {
		select {
		case job, ok := <-wp.jobCh:
			if !ok {
				return
			}
			if err := wp.client.Download(job.URL, job.Path); err != nil {
				wp.retry.Handle(wp.ctx, job, errors.New("from server"))
			}

		case <-wp.ctx.Done():
			return
		}
	}
}

func (wp *workerPool) SubmitJob(ctx context.Context, job Job) error {
	wp.mu.Lock()
	closing := wp.closing
	wp.mu.Unlock()

	if closing {
		return errors.New("worker pool is closing")
	}

	select {
	case wp.jobCh <- job:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (wp *workerPool) Close() {

	// stop accepting new jobs
	grace := 1 * time.Second
	timer := time.NewTimer(grace)
	<-timer.C

	wp.mu.Lock()
	wp.closing = true
	wp.mu.Unlock()

	// allow in-flight retries to enqueue (grace period)
	// you can tweak this duration
	timer = time.NewTimer(grace)
	<-timer.C

	// stop workers + retries
	wp.cancel()

	// wait workers to finish
	wp.wg.Wait()

	// now safe to close channel
	close(wp.jobCh)
}

func (wp *workerPool) Wait() {
	wp.wg.Wait()
}
