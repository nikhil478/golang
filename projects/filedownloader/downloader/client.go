package downloader

import (
	"io"
	"net/http"
	"os"
)

type Client struct {
	httpClient http.Client
}

func (c *Client) Download(URL string, path string) error {
	resp, err := c.httpClient.Get(URL)
	if err != nil {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	_, err = io.Copy(file, resp.Body)
	return err
}