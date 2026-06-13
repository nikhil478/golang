package jira

import (
	"net/http"
)

type Client struct {
	BaseURL string
	Token   string

	HTTP *http.Client
}

func New(baseURL, token string) *Client {
	return &Client{
		BaseURL: baseURL,
		Token:   token,
		HTTP:    &http.Client{},
	}
}
