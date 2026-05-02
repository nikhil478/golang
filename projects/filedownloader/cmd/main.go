package main

import (
	"context"

	"github.com/nikhil478/filedownloader/downloader"
)

func main() {
	wp := downloader.NewWorkerPool(context.Background(), downloader.NewClient(), 10, 10)
	wp.SubmitJob(context.Background(), downloader.Job{URL: "https://github.com/nikhil478/golang/blob/main/README.md", Path: "./output.html", MaxRetry: 2, Attempts: 0})
	wp.Close()
	wp.Wait()
}
