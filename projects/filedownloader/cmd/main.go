package main

import "github.com/nikhil478/filedownloader/downloader"

func main() {
	filedownloader := downloader.Client{}
	filedownloader.Download("https://github.com/nikhil478/golang/blob/main/README.md", "./output.html")
}
