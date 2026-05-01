package downloader

import (
	"log"
	"time"
)

type RetryHandler struct {
	pool WorkerPool
}

func NewRetryHandler(pool WorkerPool) *RetryHandler {
	return &RetryHandler{pool: pool}
}

func (r *RetryHandler) Handle(job Job, err error) {
	if job.Attempts >= job.MaxRetry {
		log.Printf("permanent failure :%v", err)
		return
	}

	job.Attempts++
	go func(j Job) {
		delay := backoff(j.Attempts)
		time.Sleep(delay)
		r.pool.SubmitJob(j)
	}(job)
}

func backoff(attempt int) time.Duration {
	return time.Duration(1<<attempt) * 100 * time.Millisecond
}
