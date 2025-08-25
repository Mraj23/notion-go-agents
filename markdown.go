package notion

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// NotionMarkdownConverter turns Notion block trees into Markdown text.
type NotionMarkdownConverter struct {
	client   *Client
	maxDepth int
	maxNodes int
}

func NewNotionMarkdownConverter(client *Client) *NotionMarkdownConverter {
	return &NotionMarkdownConverter{
		client:   client,
		maxDepth: 3,
		maxNodes: 500,
	}
}

type BlockNode struct {
	Block    map[string]any
	Children []BlockNode
}

// ConvertPageToMarkdown retrieves blocks for a page and renders them to Markdown.
func (c *NotionMarkdownConverter) ConvertPageToMarkdown(ctx context.Context, pageID string) (string, error) {
	blocks, err := c.getBlockTree(ctx, pageID, 0)
	if err != nil {
		return "", fmt.Errorf("failed to get block tree: %w", err)
	}
	var b strings.Builder
	c.renderBlocksToMarkdown(&b, blocks, 0)
	md := strings.TrimSpace(b.String())
	if md == "" {
		md = "(no textual content)"
	}
	return md, nil
}

func (c *NotionMarkdownConverter) getBlockTree(ctx context.Context, blockID string, depth int) ([]BlockNode, error) {
	if depth >= c.maxDepth {
		return nil, nil
	}
	children, err := c.getAllBlockChildren(ctx, blockID)
	if err != nil {
		return nil, err
	}
	nodes := make([]BlockNode, 0, len(children))
	for _, child := range children {
		node := BlockNode{Block: child}
		hasChildren, _ := child["has_children"].(bool)
		blockType, _ := child["type"].(string)
		var sourceID string
		if blockType == "synced_block" {
			if sb, ok := child["synced_block"].(map[string]any); ok {
				if sf, ok := sb["synced_from"].(map[string]any); ok {
					if bid, ok := sf["block_id"].(string); ok && bid != "" {
						sourceID = bid
					}
				}
			}
		}
		if sourceID == "" && hasChildren {
			if id, ok := child["id"].(string); ok && id != "" {
				sourceID = id
			}
		}
		if sourceID != "" {
			childNodes, err := c.getBlockTree(ctx, sourceID, depth+1)
			if err != nil {
				return nil, err
			}
			node.Children = childNodes
		}
		nodes = append(nodes, node)
	}
	return nodes, nil
}

func (c *NotionMarkdownConverter) getAllBlockChildren(ctx context.Context, blockID string) ([]map[string]any, error) {
	path := "/v1/blocks/" + blockID + "/children?page_size=100"
	var results []map[string]any
	cursor := ""
	for {
		p := path
		if cursor != "" {
			p = p + "&start_cursor=" + cursor
		}
		resp, err := c.client.request(ctx, http.MethodGet, p, nil)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return nil, fmt.Errorf("get children failed: status=%d body=%s", resp.StatusCode, string(body))
		}
		var blocksResp struct {
			Object     string            `json:"object"`
			Results    []json.RawMessage `json:"results"`
			NextCursor string            `json:"next_cursor"`
			HasMore    bool              `json:"has_more"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&blocksResp); err != nil {
			return nil, fmt.Errorf("failed to decode blocks: %w", err)
		}
		for _, raw := range blocksResp.Results {
			var blockMap map[string]any
			if err := json.Unmarshal(raw, &blockMap); err != nil {
				continue
			}
			results = append(results, blockMap)
			if len(results) >= c.maxNodes {
				return results, nil
			}
		}
		if !blocksResp.HasMore || blocksResp.NextCursor == "" {
			break
		}
		cursor = blocksResp.NextCursor
	}
	c.processNumberedListItems(results)
	return results, nil
}

func (c *NotionMarkdownConverter) processNumberedListItems(blocks []map[string]any) {
	numberedListIndex := 0
	for _, block := range blocks {
		blockType, _ := block["type"].(string)
		if blockType == "numbered_list_item" {
			numberedListIndex++
			if numberedItem, ok := block["numbered_list_item"].(map[string]any); ok {
				numberedItem["number"] = numberedListIndex
			}
		} else {
			numberedListIndex = 0
		}
	}
}

func (c *NotionMarkdownConverter) renderBlocksToMarkdown(b *strings.Builder, nodes []BlockNode, indent int) {
	for i, node := range nodes {
		c.renderBlockToMarkdown(b, node, indent)
		if i < len(nodes)-1 {
			b.WriteString("\n")
		}
	}
}

func (c *NotionMarkdownConverter) renderBlockToMarkdown(b *strings.Builder, node BlockNode, indent int) {
	blockType, _ := node.Block["type"].(string)
	if blockType == "column_list" || blockType == "column" || blockType == "synced_block" {
		if len(node.Children) > 0 {
			c.renderBlocksToMarkdown(b, node.Children, indent)
		}
		return
	}
	if blockType == "table" {
		c.renderTableToMarkdown(b, node, indent)
		return
	}
	line := c.blockToMarkdown(node.Block, indent)
	if strings.TrimSpace(line) != "" {
		b.WriteString(line)
	}
	if len(node.Children) > 0 {
		if strings.TrimSpace(line) != "" && !c.isListItem(blockType) {
			b.WriteString("\n")
		}
		childIndent := indent
		if c.shouldIndentChildren(blockType) {
			childIndent++
		}
		c.renderBlocksToMarkdown(b, node.Children, childIndent)
	}
}

func (c *NotionMarkdownConverter) blockToMarkdown(block map[string]any, indent int) string {
	blockType, _ := block["type"].(string)
	pad := strings.Repeat("  ", indent)
	switch blockType {
	case "heading_1":
		text := c.extractRichText(block[blockType])
		if text == "" {
			return ""
		}
		return pad + "# " + text
	case "heading_2":
		text := c.extractRichText(block[blockType])
		if text == "" {
			return ""
		}
		return pad + "## " + text
	case "heading_3":
		text := c.extractRichText(block[blockType])
		if text == "" {
			return ""
		}
		return pad + "### " + text
	case "paragraph":
		text := c.extractRichText(block[blockType])
		if text == "" {
			return ""
		}
		return pad + text
	case "bulleted_list_item":
		text := c.extractRichText(block[blockType])
		if text == "" {
			return ""
		}
		return pad + "- " + text
	case "numbered_list_item":
		text := c.extractRichText(block[blockType])
		if text == "" {
			return ""
		}
		var number int = 1
		if numberedItem, ok := block[blockType].(map[string]any); ok {
			if num, ok := numberedItem["number"].(int); ok {
				number = num
			}
		}
		return pad + strconv.Itoa(number) + ". " + text
	case "to_do":
		todoBlock, _ := block[blockType].(map[string]any)
		checked, _ := todoBlock["checked"].(bool)
		text := c.extractRichText(todoBlock)
		checkbox := "[ ]"
		if checked {
			checkbox = "[x]"
		}
		return pad + "- " + checkbox + " " + text
	case "quote":
		text := c.extractRichText(block[blockType])
		if text == "" {
			return ""
		}
		return pad + "> " + text
	case "callout":
		calloutBlock, _ := block[blockType].(map[string]any)
		text := c.extractRichText(calloutBlock)
		if text == "" {
			return ""
		}
		emoji := ""
		if icon, ok := calloutBlock["icon"].(map[string]any); ok {
			if iconType, _ := icon["type"].(string); iconType == "emoji" {
				if emojiStr, _ := icon["emoji"].(string); emojiStr != "" {
					emoji = emojiStr + " "
				}
			}
		}
		return pad + "> " + emoji + text
	case "code":
		codeBlock, _ := block[blockType].(map[string]any)
		text := c.extractRichText(codeBlock)
		if text == "" {
			return ""
		}
		language, _ := codeBlock["language"].(string)
		if language == "" {
			language = "plaintext"
		}
		return pad + "```" + language + "\n" + text + "\n" + pad + "```"
	case "divider":
		return pad + "---"
	case "bookmark", "embed", "link_preview":
		var url string
		if blockData, ok := block[blockType].(map[string]any); ok {
			url, _ = blockData["url"].(string)
		}
		if url == "" {
			return ""
		}
		return pad + "[" + blockType + "](" + url + ")"
	case "link_to_page":
		linkBlock, _ := block[blockType].(map[string]any)
		linkType, _ := linkBlock["type"].(string)
		if linkType == "page_id" {
			if pageID, _ := linkBlock["page_id"].(string); pageID != "" {
				url := "https://www.notion.so/" + pageID
				return pad + "[link](" + url + ")"
			}
		}
		return ""
	case "child_page":
		childPage, _ := block[blockType].(map[string]any)
		title, _ := childPage["title"].(string)
		if title == "" {
			return ""
		}
		return pad + "## " + title
	case "image":
		imageBlock, _ := block[blockType].(map[string]any)
		var url string
		imageType, _ := imageBlock["type"].(string)
		if imageType == "external" {
			if external, ok := imageBlock["external"].(map[string]any); ok {
				url, _ = external["url"].(string)
			}
		} else if imageType == "file" {
			if file, ok := imageBlock["file"].(map[string]any); ok {
				url, _ = file["url"].(string)
			}
		}
		alt := "image"
		if caption, ok := imageBlock["caption"].([]any); ok && len(caption) > 0 {
			var captionText strings.Builder
			for _, item := range caption {
				if itemMap, ok := item.(map[string]any); ok {
					if plainText, ok := itemMap["plain_text"].(string); ok {
						captionText.WriteString(plainText)
					}
				}
			}
			if captionText.Len() > 0 {
				alt = captionText.String()
			}
		}
		if url == "" {
			return ""
		}
		return pad + "![" + alt + "](" + url + ")"
	case "toggle":
		text := c.extractRichText(block[blockType])
		if text == "" {
			return ""
		}
		return pad + "- " + text
	case "equation":
		equationBlock, _ := block[blockType].(map[string]any)
		expression, _ := equationBlock["expression"].(string)
		if expression == "" {
			return ""
		}
		return pad + "$$\n" + expression + "\n$$"
	case "table":
		return ""
	default:
		if blockData, ok := block[blockType]; ok {
			text := c.extractRichText(blockData)
			if text != "" {
				return pad + text
			}
		}
		return ""
	}
}

func (c *NotionMarkdownConverter) extractRichText(blockData any) string {
	blockMap, ok := blockData.(map[string]any)
	if !ok {
		return ""
	}
	var richTextArray []any
	if rt, ok := blockMap["rich_text"].([]any); ok {
		richTextArray = rt
	} else if text, ok := blockMap["text"].([]any); ok {
		richTextArray = text
	}
	if len(richTextArray) == 0 {
		return ""
	}
	var result strings.Builder
	for _, item := range richTextArray {
		itemMap, ok := item.(map[string]any)
		if !ok {
			continue
		}
		if itemType, _ := itemMap["type"].(string); itemType == "equation" {
			if equation, ok := itemMap["equation"].(map[string]any); ok {
				if expr, _ := equation["expression"].(string); expr != "" {
					result.WriteString("$" + expr + "$")
				}
			}
			continue
		}
		plainText, _ := itemMap["plain_text"].(string)
		if plainText == "" {
			continue
		}
		if annotations, ok := itemMap["annotations"].(map[string]any); ok {
			plainText = c.applyAnnotations(plainText, annotations)
		}
		if href, ok := itemMap["href"].(string); ok && href != "" {
			plainText = "[" + plainText + "](" + href + ")"
		}
		result.WriteString(plainText)
	}
	return strings.TrimSpace(result.String())
}

func (c *NotionMarkdownConverter) applyAnnotations(text string, annotations map[string]any) string {
	if strings.TrimSpace(text) == "" {
		return text
	}
	leadingSpaces := ""
	trailingSpaces := ""
	if len(text) > 0 {
		for i, r := range text {
			if r != ' ' && r != '\t' && r != '\n' {
				leadingSpaces = text[:i]
				break
			}
		}
		for i := len(text) - 1; i >= 0; i-- {
			if text[i] != ' ' && text[i] != '\t' && text[i] != '\n' {
				trailingSpaces = text[i+1:]
				break
			}
		}
	}
	text = strings.TrimSpace(text)
	if text == "" {
		return leadingSpaces + trailingSpaces
	}
	if code, _ := annotations["code"].(bool); code {
		text = "`" + text + "`"
	}
	if bold, _ := annotations["bold"].(bool); bold {
		text = "**" + text + "**"
	}
	if italic, _ := annotations["italic"].(bool); italic {
		text = "_" + text + "_"
	}
	if strikethrough, _ := annotations["strikethrough"].(bool); strikethrough {
		text = "~~" + text + "~~"
	}
	if underline, _ := annotations["underline"].(bool); underline {
		text = "<u>" + text + "</u>"
	}
	return leadingSpaces + text + trailingSpaces
}

func (c *NotionMarkdownConverter) isListItem(blockType string) bool {
	return blockType == "bulleted_list_item" || blockType == "numbered_list_item" || blockType == "to_do"
}

func (c *NotionMarkdownConverter) shouldIndentChildren(blockType string) bool {
	return blockType == "bulleted_list_item" || blockType == "numbered_list_item" || blockType == "to_do" || blockType == "quote" || blockType == "callout" || blockType == "toggle"
}

func (c *NotionMarkdownConverter) renderTableToMarkdown(b *strings.Builder, node BlockNode, indent int) {
	pad := strings.Repeat("  ", indent)
	hasColumnHeader := false
	hasRowHeader := false
	numCols := 0
	if tableMeta, ok := node.Block["table"].(map[string]any); ok {
		if hch, ok := tableMeta["has_column_header"].(bool); ok {
			hasColumnHeader = hch
		}
		if hrh, ok := tableMeta["has_row_header"].(bool); ok {
			hasRowHeader = hrh
		}
		if tw, ok := tableMeta["table_width"].(float64); ok {
			numCols = int(tw)
		}
	}
	rows := make([][]string, 0, len(node.Children))
	for _, child := range node.Children {
		ctype, _ := child.Block["type"].(string)
		if ctype != "table_row" {
			continue
		}
		tr, _ := child.Block["table_row"].(map[string]any)
		cellsAny, _ := tr["cells"].([]any)
		row := make([]string, 0, len(cellsAny))
		for _, cellAny := range cellsAny {
			if rtArr, ok := cellAny.([]any); ok {
				cellText := c.extractRichText(map[string]any{"rich_text": rtArr})
				row = append(row, cellText)
			} else {
				row = append(row, "")
			}
		}
		rows = append(rows, row)
		if len(row) > numCols {
			numCols = len(row)
		}
	}
	if len(rows) == 0 || numCols == 0 {
		b.WriteString(pad + "[Table]")
		return
	}
	for i := range rows {
		if len(rows[i]) < numCols {
			missing := numCols - len(rows[i])
			for j := 0; j < missing; j++ {
				rows[i] = append(rows[i], "")
			}
		}
	}
	var header []string
	body := rows
	if hasColumnHeader && len(rows) > 0 {
		header = rows[0]
		body = rows[1:]
	} else {
		header = make([]string, numCols)
		for i := range header {
			header[i] = ""
		}
	}
	if hasRowHeader {
		for i := range body {
			if len(body[i]) > 0 {
				if body[i][0] == "" {
					body[i][0] = "** **"
				} else {
					body[i][0] = "**" + body[i][0] + "**"
				}
			}
		}
	}
	b.WriteString(pad + "| " + strings.Join(header, " | ") + " |\n")
	sepCells := make([]string, numCols)
	for i := 0; i < numCols; i++ {
		sepCells[i] = "---"
	}
	b.WriteString(pad + "| " + strings.Join(sepCells, " | ") + " |")
	for _, r := range body {
		b.WriteString("\n" + pad + "| " + strings.Join(r, " | ") + " |")
	}
}
