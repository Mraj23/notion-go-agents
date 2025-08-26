Here’s a clean, merged README with the latest from both sides—no conflict markers, unified sections, and consistent naming.

# Notion KB for Go Agents

A lightweight Go package for using Notion as a knowledge base in LLM-powered agents. It wraps Notion’s REST API with a small, composable client, provides helpers for querying databases and pages, and includes a Markdown converter for Notion blocks.

## Features

* **Typed client**: Minimal `Client` for Notion REST calls.
* **Database querying**: Fetch pages from a Notion database **with pagination support**.
* **Workspace search**: Search Notion and filter to pages.
* **Page retrieval**: Fetch page metadata and properties.
* **Markdown conversion**: Convert Notion page blocks into readable Markdown.
* **Utilities**: `ExtractNotionTitle`, `SelectPrintableProperties`.
* **Convenience helpers**: `FindNotionPage`, `SearchNotionDatabase` (aka `SearchNotionDB`) and `GetPageContent` to wrap common tasks that only need an API key.

## Installation

```bash
go get github.com/openai/notion-go-agents
```

## Quick Start

### One-liner: find a page and print Markdown

```go
package main

import (
    "context"
    "fmt"
    "log"

    notion "github.com/openai/notion-go-agents"
)

func main() {
    apiKey := "YOUR_NOTION_API_KEY"
    ctx := context.Background()

    // Find the first page matching a query and print its Markdown
    page, err := notion.FindNotionPage(ctx, apiKey, "roadmap")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(page.Title)
    fmt.Println(page.Markdown)
}
```

### Client-based usage (search workspace)

```go
package main

import (
    "context"
    "fmt"
    "os"
    "time"

    notion "github.com/openai/notion-go-agents"
)

func main() {
    ctx := context.Background()
    client := notion.NewClient(
        os.Getenv("NOTION_API_KEY"),
        os.Getenv("NOTION_VERSION"),          // optional; defaults to 2022-06-28
        notion.WithTimeout(10*time.Second),   // optional
        // notion.WithHTTPClient(customHTTPClient),
    )

    pages, _ := notion.SearchWorkspace(ctx, client, notion.NotionSearchRequest{Query: "docs"}, 1)
    for _, pg := range pages {
        fmt.Println(pg.Title)
        fmt.Println(pg.Properties)
        fmt.Println(pg.Markdown)
    }
}
```

## Folder Structure

```
pkg/notion/
  client.go   — Notion Client and HTTP request wrapper
  types.go    — Request/response and model types (search, database, page)
  api.go      — High-level API methods: SearchPages, GetPage, QueryDatabase
  helpers.go  — Higher-level helpers: SearchWorkspace, SearchNotionDatabase (SearchNotionDB), GetPageContent
  extract.go  — Helpers: ExtractNotionTitle, SelectPrintableProperties
  markdown.go — NotionMarkdownConverter to render blocks as Markdown
```

## Requirements

* **Go**: 1.21+ (recommended)
* **Notion Integration**:

  * Create a Notion integration and copy its API key.
  * Share the relevant pages/databases with the integration.
* **Environment variables**:

  * `NOTION_API_KEY` (required)
  * `NOTION_VERSION` (optional; defaults to `2022-06-28`)

## Examples

* `examples/database/main.go` — Search a Notion database with `SearchNotionDatabase`.
* More examples are available in the [`examples/`](./examples) directory.

## Contributing

See [CONTRIBUTING.md](./CONTRIBUTING.md) for guidelines.

## License

MIT
