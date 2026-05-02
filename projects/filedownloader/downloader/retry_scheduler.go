package downloader

import (
	"container/heap"
	"context"
	"sync"
	"time"
)

type JobSubmitter interface {
	SubmitJob(ctx context.Context, job Job) error
}

type RetryScheduler struct {
	submitter JobSubmitter

	mu   sync.Mutex
	heap jobHeap
	wake chan struct{}

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func NewRetryScheduler(parent context.Context, submitter JobSubmitter) *RetryScheduler {
	ctx, cancel := context.WithCancel(parent)

	rs := &RetryScheduler{
		submitter: submitter,
		wake:      make(chan struct{}, 1),
		ctx:       ctx,
		cancel:    cancel,
	}

	heap.Init(&rs.heap)

	rs.wg.Add(1)
	go rs.loop()

	return rs
}

func (rs *RetryScheduler) Schedule(job Job, delay time.Duration) {
	runAt := time.Now().Add(delay)

	rs.mu.Lock()
	heap.Push(&rs.heap, &scheduledJob{
		job:   job,
		runAt: runAt,
	})
	rs.mu.Unlock()

	// notify loop
	select {
	case rs.wake <- struct{}{}:
	default:
	}
}

func (rs *RetryScheduler) loop() {
	defer rs.wg.Done()

	for {
		rs.mu.Lock()

		if len(rs.heap) == 0 {
			rs.mu.Unlock()

			select {
			case <-rs.wake:
				continue
			case <-rs.ctx.Done():
				return
			}
		}

		next := rs.heap[0]
		wait := time.Until(next.runAt)

		rs.mu.Unlock()

		if wait > 0 {
			select {
			case <-time.After(wait):
			case <-rs.wake:
				continue
			case <-rs.ctx.Done():
				return
			}
		}

		// pop and execute
		rs.mu.Lock()
		if len(rs.heap) == 0 {
			rs.mu.Unlock()
			continue
		}

		next = heap.Pop(&rs.heap).(*scheduledJob)
		rs.mu.Unlock()

		// submit (non-blocking responsibility of submitter)
		_ = rs.submitter.SubmitJob(rs.ctx, next.job)
	}
}

func (rs *RetryScheduler) Close() {
	rs.cancel()
	rs.wg.Wait()
}
