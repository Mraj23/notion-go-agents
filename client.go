package notion

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	httpClient    *http.Client
	apiKey        string
	notionVersion string
}

func NewClient(apiKey, notionVersion string) *Client {
	if notionVersion == "" {
		notionVersion = "2022-06-28"
	}
	return &Client{
		httpClient:    &http.Client{Timeout: 30 * time.Second},
		apiKey:        apiKey,
		notionVersion: notionVersion,
	}
}

func (c *Client) request(ctx context.Context, method, path string, body io.Reader) (*http.Response, error) {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	url := "https://api.notion.com" + path
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Notion-Version", c.notionVersion)
	if method == http.MethodPost || method == http.MethodPatch {
		req.Header.Set("Content-Type", "application/json")
	}
	return c.httpClient.Do(req)
}
