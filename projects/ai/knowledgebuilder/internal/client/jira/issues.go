package jira

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Issue struct {
	Key string `json:"key"`
}

func (c *Client) GetIssue(
	issueKey string,
) (*Issue, error) {

	url := fmt.Sprintf(
		"%s/rest/api/3/issue/%s",
		c.BaseURL,
		issueKey,
	)

	req, err := http.NewRequest(
		http.MethodGet,
		url,
		nil,
	)

	if err != nil {
		return nil, err
	}

	req.Header.Set(
		"Authorization",
		"Bearer "+c.Token,
	)

	resp, err := c.HTTP.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var issue Issue

	err = json.NewDecoder(
		resp.Body,
	).Decode(&issue)

	if err != nil {
		return nil, err
	}

	return &issue, nil
}
