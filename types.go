package notion

import "encoding/json"

type NotionSearchRequest struct {
	Query  string           `json:"query,omitempty"`
	Filter *NotionObjFilter `json:"filter,omitempty"`
	Sort   *NotionSort      `json:"sort,omitempty"`
}

type NotionObjFilter struct {
	Value    string `json:"value,omitempty"`
	Property string `json:"property,omitempty"`
}

type NotionSort struct {
	Direction string `json:"direction,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
}

type NotionSearchResponse struct {
	Object     string            `json:"object"`
	Results    []json.RawMessage `json:"results"`
	HasMore    bool              `json:"has_more"`
	NextCursor string            `json:"next_cursor"`
}

type NotionPageRef struct {
	Object string `json:"object"`
	ID     string `json:"id"`
	URL    string `json:"url"`
}

type NotionPage struct {
	NotionPageRef
	CreatedTime    string         `json:"created_time"`
	LastEditedTime string         `json:"last_edited_time"`
	Parent         map[string]any `json:"parent"`
	Archived       bool           `json:"archived"`
	Properties     map[string]any `json:"properties"`
	PublicURL      string         `json:"public_url"`
}

type NotionDatabaseQueryRequest struct {
	Filter      map[string]any `json:"filter,omitempty"`
	Sorts       []NotionSort   `json:"sorts,omitempty"`
	StartCursor string         `json:"start_cursor,omitempty"`
	PageSize    int            `json:"page_size,omitempty"`
}

type NotionDatabaseQueryResponse struct {
	Object     string            `json:"object"`
	Results    []json.RawMessage `json:"results"`
	HasMore    bool              `json:"has_more"`
	NextCursor string            `json:"next_cursor"`
}
