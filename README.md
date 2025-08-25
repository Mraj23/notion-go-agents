Notion KB for Go Agents
A lightweight Go package for using Notion as a knowledge base in LLM-powered agents. It wraps Notion’s REST API with a small, composable client, provides helpers for querying databases and pages, and includes a Markdown converter for Notion blocks.

Features
Typed client: Minimal Client for Notion REST calls.
Database querying: Fetch pages from a Notion database (with pagination support).
Workspace search: Search Notion and filter to pages.
Page retrieval: Fetch page metadata and properties.
Markdown conversion: Convert Notion page blocks into readable Markdown.
Utilities: Extract page titles and select printable properties.

Example
```
client := notion.NewClient(os.Getenv("NOTION_API_KEY"), "")
pages, _ := notion.SearchWorkspace(ctx, client, notion.NotionSearchRequest{Query: "docs"}, 1)
for _, pg := range pages {
    fmt.Println(pg.Title)
    fmt.Println(pg.Properties)
    fmt.Println(pg.Markdown)
}
```

Folder structure
pkg/notion/
client.go — Notion Client and HTTP request wrapper
types.go — Request/response and model types (search, database, page)
api.go — High-level API methods: SearchPages, GetPage, QueryDatabase
helpers.go — Higher level helpers: SearchWorkspace, SearchNotionDatabase, GetPageContent
extract.go — Helpers: ExtractNotionTitle, SelectPrintableProperties
markdown.go — NotionMarkdownConverter to render blocks as Markdown

Requirements
Go: 1.21+ (recommended)
Notion Integration:
Create a Notion integration and copy its API key.
Share the relevant pages/databases with the integration.
Environment variables:
NOTION_API_KEY (required)
NOTION_VERSION (optional, defaults to 2022-06-28)

More examples
examples/database/main.go — Search a Notion database with SearchNotionDatabase.

Quick Start

```go
client := notion.NewClient(
    os.Getenv("NOTION_API_KEY"),
    os.Getenv("NOTION_VERSION"),
    notion.WithTimeout(10*time.Second),
    // notion.WithHTTPClient(customHTTPClient),
)
```