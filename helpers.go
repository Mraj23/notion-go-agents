package notion

import (
	"context"
	"encoding/json"
	"fmt"
)

// PageContent contains basic page information returned by helpers.
type PageContent struct {
	ID         string
	Title      string
	Markdown   string
	URL        string
	Properties map[string]any
}

// GetPageContent retrieves a Notion page by ID and converts it to Markdown.
func GetPageContent(ctx context.Context, client *Client, pageID string) (*PageContent, error) {
	pg, err := client.GetPage(ctx, pageID)
	if err != nil {
		return nil, err
	}
	conv := NewNotionMarkdownConverter(client)
	md, err := conv.ConvertPageToMarkdown(ctx, pageID)
	if err != nil {
		return nil, fmt.Errorf("failed to convert page to markdown: %w", err)
	}
	title := ExtractNotionTitle(pg.Properties)
	props := SelectPrintableProperties(pg.Properties)

	url := pg.PublicURL
	if url == "" {
		url = pg.URL
	}

	return &PageContent{
		ID:         pg.ID,
		Title:      title,
		Markdown:   md,
		URL:        url,
		Properties: props,
	}, nil
}

// SearchWorkspace searches the workspace and returns up to limit page results.
func SearchWorkspace(ctx context.Context, client *Client, req NotionSearchRequest, limit int) ([]PageContent, error) {
	ids, err := client.SearchPages(ctx, req, limit)
	if err != nil {
		return nil, err
	}
	out := make([]PageContent, 0, len(ids))
	for _, id := range ids {
		pc, err := GetPageContent(ctx, client, id)
		if err != nil {
			continue
		}
		out = append(out, *pc)
		if limit > 0 && len(out) >= limit {
			break
		}
	}
	return out, nil
}

// SearchNotionDatabase queries a database (paginated) and returns up to limit page results.
func SearchNotionDatabase(ctx context.Context, client *Client, databaseID string, req NotionDatabaseQueryRequest, limit int) ([]PageContent, error) {
	// Use large page size and paginate internally; enforce overall limit manually.
	req.PageSize = 100

	var out []PageContent
	cursor := ""
	for {
		req.StartCursor = cursor
		resp, err := client.QueryDatabase(ctx, databaseID, req)
		if err != nil {
			return nil, err
		}

		for _, raw := range resp.Results {
			var pg NotionPage
			if err := json.Unmarshal(raw, &pg); err != nil {
				continue
			}
			pc, err := GetPageContent(ctx, client, pg.ID)
			if err != nil {
				continue
			}
			out = append(out, *pc)
			if limit > 0 && len(out) >= limit {
				return out[:limit], nil
			}
		}

		if !resp.HasMore || resp.NextCursor == "" {
			break
		}
		cursor = resp.NextCursor
	}
	return out, nil
}

// FindPageByQuery is a small convenience wrapper that returns the first matching page.
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
