package main

import (
	"context"

	"github.com/nikhil478/filedownloader/downloader"
)

func main() {
	wp := downloader.NewWorkerPool(downloader.NewClient())
	wp.Run(context.Background(), 10, 10)
	wp.SubmitJob(downloader.Job{URL: "https://github.com/nikhil478/golang/blob/main/README.md", Path: "./output.html"})
	wp.Close()
	wp.Wait()
}
