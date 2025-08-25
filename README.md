# Notion KB for Go Agents

A lightweight Go package for using Notion as a knowledge base in LLM-powered agents. It wraps Notion’s REST API with a small, composable client, provides helpers for querying databases and pages, and includes a Markdown converter for Notion blocks.

## Features
- **Typed client**: Minimal `Client` for Notion REST calls.
- **Database querying**: Fetch pages from a Notion database.
- **Workspace search**: Search Notion and filter to pages.
- **Page retrieval**: Fetch page metadata and properties.
- **Markdown conversion**: Convert Notion page blocks into readable Markdown.
- **Utilities**: Extract page titles and select printable properties.
- **Convenience helpers**: `FindNotionPage` and `SearchNotionDB` wrap common
  tasks and require only an API key.

## Installation
```bash
go get github.com/openai/notion-go-agents
```

## Quick Start
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

More examples can be found in the [`examples/`](./examples) directory.

## Contributing
See [CONTRIBUTING.md](./CONTRIBUTING.md) for guidelines.

## License
MIT
