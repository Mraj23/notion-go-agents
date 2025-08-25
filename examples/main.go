package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"notion"
)

func main() {
	apiKey := os.Getenv("NOTION_API_KEY")
	if apiKey == "" {
		log.Fatal("NOTION_API_KEY is required")
	}
	client := notion.NewClient(apiKey, "")
	ctx := context.Background()
	pages, err := notion.SearchWorkspace(ctx, client, notion.NotionSearchRequest{Query: "docs"}, 1)
	if err != nil {
		log.Fatal(err)
	}
	for _, pg := range pages {
		fmt.Println("Title:", pg.Title)
		fmt.Println("Properties:", pg.Properties)
		fmt.Println("Content:", pg.Markdown)
	}
}
