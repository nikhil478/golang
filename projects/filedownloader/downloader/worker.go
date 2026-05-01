package downloader

import (
	"context"
	"errors"
	"sync"
)

type WorkerPool interface {
	SubmitJob(job Job)
	Close()
	Wait()
}

type workerPool struct {
	client *Client
	wg     sync.WaitGroup
	jobCh  chan Job

	retry *RetryHandler
}

func NewWorkerPool(ctx context.Context, client *Client, workers int, buffer int) WorkerPool {
	wp := &workerPool{client: client}
	wp.jobCh = make(chan Job, buffer)
	wp.retry = NewRetryHandler(wp)
	for range workers {
		wp.wg.Go(func() {
			wp.worker(ctx, wp.jobCh)
		})
	}
	return wp
}

func (wp *workerPool) worker(ctx context.Context, jobCh chan Job) {
	for {
		select {
		case job, ok := <-jobCh:
			if !ok {
				return
			}
			if err := wp.client.Download(job.URL, job.Path); err != nil {
				wp.retry.Handle(job, errors.New("error while downloading"))
			}
		case <-ctx.Done():
			return
		}
	}
}

func (wp *workerPool) SubmitJob(job Job) {
	wp.jobCh <- job
}

func (wp *workerPool) Close() {
	close(wp.jobCh)
}

func (wp *workerPool) Wait() {
	wp.wg.Wait()
}
