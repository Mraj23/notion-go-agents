package notion

import "encoding/json"

// NotionSearchRequest describes parameters for the search API.
type NotionSearchRequest struct {
	Query  string           `json:"query,omitempty"`
	Filter *NotionObjFilter `json:"filter,omitempty"`
	Sort   *NotionSort      `json:"sort,omitempty"`
}

// NotionObjFilter filters search results by object type.
type NotionObjFilter struct {
	Value    string `json:"value,omitempty"`
	Property string `json:"property,omitempty"`
}

// NotionSort specifies sort options for search and database queries.
type NotionSort struct {
	Direction string `json:"direction,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
}

// NotionSearchResponse is the payload returned from the search API.
type NotionSearchResponse struct {
	Object     string            `json:"object"`
	Results    []json.RawMessage `json:"results"`
	HasMore    bool              `json:"has_more"`
	NextCursor string            `json:"next_cursor"`
}

// NotionPageRef is a lightweight reference to a page.
type NotionPageRef struct {
	Object string `json:"object"`
	ID     string `json:"id"`
	URL    string `json:"url"`
}

// NotionPage represents a page object with properties and metadata.
type NotionPage struct {
	NotionPageRef
	CreatedTime    string         `json:"created_time"`
	LastEditedTime string         `json:"last_edited_time"`
	Parent         map[string]any `json:"parent"`
	Archived       bool           `json:"archived"`
	Properties     map[string]any `json:"properties"`
	PublicURL      string         `json:"public_url"`
}

// NotionDatabaseQueryRequest describes a database query.
type NotionDatabaseQueryRequest struct {
	Filter      map[string]any `json:"filter,omitempty"`
	Sorts       []NotionSort   `json:"sorts,omitempty"`
	StartCursor string         `json:"start_cursor,omitempty"`
	PageSize    int            `json:"page_size,omitempty"`
}

// NotionDatabaseQueryResponse holds results from a database query.
type NotionDatabaseQueryResponse struct {
	Object     string            `json:"object"`
	Results    []json.RawMessage `json:"results"`
	HasMore    bool              `json:"has_more"`
	NextCursor string            `json:"next_cursor"`
}
