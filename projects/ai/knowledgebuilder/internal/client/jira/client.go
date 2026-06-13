package jira

import (
	"net/http"
)

type Client struct {
	BaseURL string
	Email   string
	Token   string

	HTTP *http.Client
}

func New(
	baseURL string,
	email string,
	token string,
) *Client {

	return &Client{
		BaseURL: baseURL,
		Email:   email,
		Token:   token,
		HTTP:    &http.Client{},
	}
}
