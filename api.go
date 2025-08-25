package notion

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// SearchPages queries the Notion search API and returns page IDs.
// The limit parameter restricts the number of page IDs returned (0 means no limit).
func (c *Client) SearchPages(ctx context.Context, req NotionSearchRequest, limit int) ([]string, error) {
	bts, _ := json.Marshal(req)
	resp, err := c.request(ctx, http.MethodPost, "/v1/search", bytes.NewReader(bts))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("notion search failed: status=%d body=%s", resp.StatusCode, string(body))
	}
	var sr NotionSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&sr); err != nil {
		return nil, fmt.Errorf("failed to decode search response: %w", err)
	}
	ids := make([]string, 0, limit)
	for _, raw := range sr.Results {
		var ref NotionPageRef
		if err := json.Unmarshal(raw, &ref); err != nil {
			continue
		}
		if ref.Object != "page" {
			continue
		}
		ids = append(ids, ref.ID)
		if limit > 0 && len(ids) >= limit {
			break
		}
	}
	return ids, nil
}

// GetPage fetches metadata for a Notion page by its ID.
func (c *Client) GetPage(ctx context.Context, pageID string) (*NotionPage, error) {
	path := "/v1/pages/" + pageID
	resp, err := c.request(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("get page failed: status=%d body=%s", resp.StatusCode, string(body))
	}
	var pg NotionPage
	if err := json.NewDecoder(resp.Body).Decode(&pg); err != nil {
		return nil, fmt.Errorf("failed to decode page: %w", err)
	}
	return &pg, nil
}

// QueryDatabase runs a query against a Notion database and returns the raw response.
func (c *Client) QueryDatabase(ctx context.Context, databaseID string, req NotionDatabaseQueryRequest) (*NotionDatabaseQueryResponse, error) {
	path := "/v1/databases/" + databaseID + "/query"
	bts, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query request: %w", err)
	}
	resp, err := c.request(ctx, http.MethodPost, path, bytes.NewReader(bts))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("database query failed: status=%d body=%s", resp.StatusCode, string(body))
	}
	var out NotionDatabaseQueryResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("failed to decode query response: %w", err)
	}
	return &out, nil
}
