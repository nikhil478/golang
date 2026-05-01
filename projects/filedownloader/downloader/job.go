package downloader

type Job struct {
	URL  string
	Path string

	Attempts int
	MaxRetry int
}
