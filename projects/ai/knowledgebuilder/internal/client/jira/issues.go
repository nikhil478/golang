package jira

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type Issue struct {
	Key    string      `json:"key"`
	Fields IssueFields `json:"fields"`
}

type IssueFields struct {
	Summary     string       `json:"summary"`
	Description any          `json:"description"`
	Attachments []Attachment `json:"attachment"`
}

type Attachment struct {
	ID       string `json:"id"`
	FileName string `json:"filename"`
	MimeType string `json:"mimeType"`
	Content  string `json:"content"`
	Size     int64  `json:"size"`
}

func (c *Client) GetIssue(
	issueKey string,
) (*Issue, error) {

	url := fmt.Sprintf(
		"%s/rest/api/3/issue/%s?fields=attachment,summary,description",
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

	req.SetBasicAuth(
		c.Email,
		c.Token,
	)

	req.Header.Set(
		"Accept",
		"application/json",
	)

	resp, err := c.HTTP.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(body))

	var issue Issue

	err = json.Unmarshal(body, &issue)
	if err != nil {
		return nil, err
	}

	return &issue, nil
}

func (c *Client) DownloadAttachments(
	issue *Issue,
	outputDir string,
) error {

	if err := os.MkdirAll(
		outputDir,
		0755,
	); err != nil {
		return err
	}

	fmt.Printf(
		"found %d attachments\n",
		len(issue.Fields.Attachments),
	)

	for _, attachment := range issue.Fields.Attachments {

		fmt.Printf(
			"downloading %s\n",
			attachment.FileName,
		)

		err := c.downloadFile(
			attachment.Content,
			filepath.Join(
				outputDir,
				attachment.FileName,
			),
		)

		if err != nil {

			return fmt.Errorf(
				"failed downloading %s: %w",
				attachment.FileName,
				err,
			)
		}
	}

	return nil
}

func (c *Client) downloadFile(
	url string,
	outputPath string,
) error {

	req, err := http.NewRequest(
		http.MethodGet,
		url,
		nil,
	)

	if err != nil {
		return err
	}

	req.SetBasicAuth(
    c.Email,
    c.Token,
)

	resp, err := c.HTTP.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {

		return fmt.Errorf(
			"unexpected status code %d",
			resp.StatusCode,
		)
	}

	file, err := os.Create(
		outputPath,
	)

	if err != nil {
		return err
	}

	defer file.Close()

	_, err = io.Copy(
		file,
		resp.Body,
	)

	return err
}
