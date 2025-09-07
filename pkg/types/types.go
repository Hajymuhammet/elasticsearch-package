package types

type IndexResponse struct {
	Index   string `json:"_index,omitempty"`
	ID      string `json:"_id,omitempty"`
	Version int64  `json:"_version,omitempty"`
	Result  string `json:"result,omitempty"`
}

type BulkResponse struct {
	Took   int  `json:"took"`
	Errors bool `json:"errors"`
	// Items omitted for brevity
}

type SearchResponse[T any] struct {
	Hits  []T
	Total int64
}

type SearchHit[T any] struct {
	Source T `json:"_source"`
}
