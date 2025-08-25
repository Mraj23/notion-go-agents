package main

import (
	"context"
	"fmt"
	"log"
	"os"

	notion "github.com/openai/notion-go-agents"
)

func main() {
	apiKey := os.Getenv("NOTION_API_KEY")
	if apiKey == "" {
		log.Fatal("NOTION_API_KEY not set")
	}
	client := notion.NewClient(apiKey, "")
	ctx := context.Background()
	page, err := notion.FindPageByQuery(ctx, client, "welcome")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(page.Title)
	fmt.Println(page.Markdown)
}
