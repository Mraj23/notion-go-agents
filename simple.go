package notion

import "context"

// FindNotionPage searches workspace pages for the given query using a new client
// created with the provided apiKey. It returns the first matching page as
// Markdown-formatted content.
func FindNotionPage(ctx context.Context, apiKey, query string) (*PageContent, error) {
	client := NewClient(apiKey, "")
	return FindPageByQuery(ctx, client, query)
}

// SearchNotionDB queries a Notion database for pages whose title contains the
// query. The client is constructed from apiKey. Limit controls the maximum number
// of pages returned (0 means no limit).
func SearchNotionDB(ctx context.Context, apiKey, databaseID, query string, limit int) ([]PageContent, error) {
	client := NewClient(apiKey, "")
	return SearchNotionDatabase(ctx, client, databaseID, query, limit)
}
