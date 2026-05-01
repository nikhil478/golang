package downloader

import (
	"context"
	"log"
	"sync"
)

type Job struct {
	URL  string
	Path string
}

type WorkerPool struct {
	client *Client
	wg     sync.WaitGroup
	jobCh  chan Job
}

func NewWorkerPool(client *Client) *WorkerPool {
	return &WorkerPool{client: client}
}

func (wp *WorkerPool) Run(ctx context.Context, workers int, buffer int) {
	wp.jobCh = make(chan Job, buffer)
	for range workers {
		wp.wg.Go(func() {
			wp.worker(ctx, wp.jobCh)
		})
	}
}

func (wp *WorkerPool) worker(ctx context.Context, jobCh chan Job) {
	for {
		select {
		case job, ok := <-jobCh:
			if !ok {
				return
			}
			if err := wp.client.Download(job.URL, job.Path); err != nil {
				log.Printf("download failed: %v", err)
			}
		case <-ctx.Done():
			return
		}
	}
}

func (wp *WorkerPool) SubmitJob(job Job) {
	wp.jobCh <- job
}

func (wp *WorkerPool) Close() {
	close(wp.jobCh)
}

func (wp *WorkerPool) Wait() {
	wp.wg.Wait()
}
