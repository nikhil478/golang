package downloader

import (
	"time"
)

type RetryHandler struct {
	scheduler *RetryScheduler
}

func NewRetryHandler(s *RetryScheduler) *RetryHandler {
	return &RetryHandler{scheduler: s}
}

func (r *RetryHandler) Handle(job Job, err error) {
	if job.Attempts >= job.MaxRetry {
		return
	}

	job.Attempts++
	delay := backoff(job.Attempts)

	r.scheduler.Schedule(job, delay)
}

func backoff(attempt int) time.Duration {
	return time.Duration(1<<attempt) * 100 * time.Millisecond
}
