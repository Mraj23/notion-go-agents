package notion

import (
	"context"
	"encoding/json"
	"fmt"
)

// PageContent contains basic page information returned by helpers.
type PageContent struct {
	ID       string
	Title    string
	Markdown string
}

// GetPageContent retrieves a Notion page by ID and converts it to Markdown.
func GetPageContent(ctx context.Context, client *Client, pageID string) (*PageContent, error) {
	pg, err := client.GetPage(ctx, pageID)
	if err != nil {
		return nil, err
	}
	converter := NewNotionMarkdownConverter(client)
	md, err := converter.ConvertPageToMarkdown(ctx, pageID)
	if err != nil {
		return nil, err
	}
	title := ExtractNotionTitle(pg.Properties)
	return &PageContent{ID: pg.ID, Title: title, Markdown: md}, nil
}

// FindPageByQuery searches workspace pages and returns the first match.
func FindPageByQuery(ctx context.Context, client *Client, query string) (*PageContent, error) {
	ids, err := client.SearchPages(ctx, NotionSearchRequest{Query: query}, 1)
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return nil, fmt.Errorf("no pages found")
	}
	return GetPageContent(ctx, client, ids[0])
}

// SearchNotionDatabase queries a database for pages whose title contains the query.
// Limit controls the maximum number of pages returned (0 means no limit).
func SearchNotionDatabase(ctx context.Context, client *Client, databaseID, query string, limit int) ([]PageContent, error) {
	filter := map[string]any{
		"or": []map[string]any{
			{"property": "Name", "title": map[string]any{"contains": query}},
			{"property": "Title", "title": map[string]any{"contains": query}},
			{"property": "name", "title": map[string]any{"contains": query}},
			{"property": "title", "title": map[string]any{"contains": query}},
		},
	}
	req := NotionDatabaseQueryRequest{Filter: filter}
	if limit > 0 {
		req.PageSize = limit
	}
	resp, err := client.QueryDatabase(ctx, databaseID, req)
	if err != nil {
		return nil, err
	}
	out := make([]PageContent, 0, len(resp.Results))
	for _, raw := range resp.Results {
		var pg NotionPage
		if err := json.Unmarshal(raw, &pg); err != nil {
			continue
		}
		content, err := GetPageContent(ctx, client, pg.ID)
		if err != nil {
			continue
		}
		out = append(out, *content)
	}
	return out, nil
}

// SearchWorkspace searches the workspace for pages matching the query and returns up to limit results.
func SearchWorkspace(ctx context.Context, client *Client, query string, limit int) ([]PageContent, error) {
	ids, err := client.SearchPages(ctx, NotionSearchRequest{Query: query}, limit)
	if err != nil {
		return nil, err
	}
	out := make([]PageContent, 0, len(ids))
	for _, id := range ids {
		content, err := GetPageContent(ctx, client, id)
		if err != nil {
			continue
		}
		out = append(out, *content)
	}
	return out, nil
}
