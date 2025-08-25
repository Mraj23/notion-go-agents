package notion

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type ClientOption func(*Client)

type Client struct {
	httpClient    *http.Client
	apiKey        string
	notionVersion string
	timeout       time.Duration
}

func NewClient(apiKey, notionVersion string, opts ...ClientOption) *Client {
	if notionVersion == "" {
		notionVersion = "2022-06-28"
	}
	c := &Client{
		apiKey:        apiKey,
		notionVersion: notionVersion,
		timeout:       30 * time.Second,
	}
	for _, opt := range opts {
		opt(c)
	}
	if c.httpClient == nil {
		c.httpClient = &http.Client{}
	}
	if c.timeout > 0 {
		c.httpClient.Timeout = c.timeout
	} else if c.httpClient.Timeout == 0 {
		c.httpClient.Timeout = 30 * time.Second
	}
	return c
}

func WithHTTPClient(h *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = h
		if c.timeout > 0 && c.httpClient != nil {
			c.httpClient.Timeout = c.timeout
		}
	}
}

func WithTimeout(d time.Duration) ClientOption {
	return func(c *Client) {
		c.timeout = d
		if c.httpClient != nil {
			c.httpClient.Timeout = d
		}
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
