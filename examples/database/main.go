package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	notion "github.com/your_org/notion-go-agents"
)

func main() {
	apiKey := os.Getenv("NOTION_API_KEY")
	databaseID := os.Getenv("DATABASE_ID")
	if apiKey == "" || databaseID == "" {
		log.Fatal("NOTION_API_KEY and DATABASE_ID must be set")
	}
	client := notion.NewClient(apiKey, os.Getenv("NOTION_VERSION"))
	ctx := context.Background()
	pages, err := SearchNotionDatabase(ctx, client, databaseID)
	if err != nil {
		log.Fatalf("search failed: %v", err)
	}
	for _, pg := range pages {
		title := notion.ExtractNotionTitle(pg.Properties)
		props := notion.SelectPrintableProperties(pg.Properties)
		bts, _ := json.MarshalIndent(props, "", "  ")
		fmt.Printf("%s\n%s\n\n", title, bts)
	}
}

func SearchNotionDatabase(ctx context.Context, c *notion.Client, databaseID string) ([]*notion.NotionPage, error) {
	resp, err := c.QueryDatabase(ctx, databaseID, notion.NotionDatabaseQueryRequest{})
	if err != nil {
		return nil, err
	}
	pages := make([]*notion.NotionPage, 0, len(resp.Results))
	for _, raw := range resp.Results {
		var pg notion.NotionPage
		if err := json.Unmarshal(raw, &pg); err != nil {
			continue
		}
		pages = append(pages, &pg)
	}
	return pages, nil
}
