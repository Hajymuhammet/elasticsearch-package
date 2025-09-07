package types

// IndexResponse represents the response from an index operation
type IndexResponse struct {
	Index   string `json:"_index,omitempty"`
	ID      string `json:"_id,omitempty"`
	Version int64  `json:"_version,omitempty"`
	Result  string `json:"result,omitempty"`
}

// BulkResponse represents the response from a bulk operation
type BulkResponse struct {
	Took   int  `json:"took"`
	Errors bool `json:"errors"`
	// TODO: Add Items if needed for detailed per-document result
}

// SearchResponse represents the response from a search query
type SearchResponse[T any] struct {
	Took     int            `json:"took"`
	TimedOut bool           `json:"timed_out"`
	Hits     SearchHits[T]  `json:"hits"`
	Shards   map[string]any `json:"_shards,omitempty"`
}

// SearchHits wraps hit metadata and documents
type SearchHits[T any] struct {
	Total    TotalCount     `json:"total"`
	MaxScore float64        `json:"max_score"`
	Hits     []SearchHit[T] `json:"hits"`
}

// TotalCount stores the total number of hits
type TotalCount struct {
	Value    int64  `json:"value"`
	Relation string `json:"relation"`
}

// SearchHit represents a single search hit (document)
type SearchHit[T any] struct {
	Index  string  `json:"_index"`
	ID     string  `json:"_id"`
	Score  float64 `json:"_score"`
	Source T       `json:"_source"`
}
