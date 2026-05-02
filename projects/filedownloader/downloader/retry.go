package downloader

import (
	"context"
	"fmt"
	"time"
)

type JobSubmitter interface {
	SubmitJob(ctx context.Context, job Job) error
}

type RetryHandler struct {
	submitter JobSubmitter
}

func NewRetryHandler(submitter JobSubmitter) *RetryHandler {
	return &RetryHandler{
		submitter: submitter,
	}
}

func (r *RetryHandler) Handle(ctx context.Context, job Job, err error) {
	if job.Attempts >= job.MaxRetry {
		return
	}

	job.Attempts++
	delay := backoff(job.Attempts)

	go func() {
		fmt.Println("retry scheduled")

		select {
		case <-time.After(delay):
			fmt.Println("retry firing")
			err := r.submitter.SubmitJob(ctx, job)
			fmt.Println("retry submit result:", err)
		case <-ctx.Done():
			fmt.Println("retry cancelled")
		}
	}()
}

func backoff(attempt int) time.Duration {
	return time.Duration(1<<attempt) * 100 * time.Millisecond
}
