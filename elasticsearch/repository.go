package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	es "github.com/elastic/go-elasticsearch/v8"

	appErrors "github.com/salamsites/elasticsearch-package/pkg/errors"
	"github.com/salamsites/elasticsearch-package/pkg/types"
)

// Repository implements repository.IndexRepository[T] using go-elasticsearch
type Repository[T any] struct {
	client  *es.Client
	httpCli *http.Client // optional, exposed if custom HTTP calls needed
}

func NewRepository[T any](client *es.Client) *Repository[T] {
	return &Repository[T]{client: client}
}

// EnsureIndex checks existence and creates index with mapping if missing
func (r *Repository[T]) EnsureIndex(ctx context.Context, index string, mapping map[string]any) error {
	res, err := r.client.Indices.Exists([]string{index}, r.client.Indices.Exists.WithContext(ctx))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode == 200 {
		return nil
	}

	body, err := json.Marshal(mapping)
	if err != nil {
		return err
	}

	createRes, err := r.client.Indices.Create(
		index,
		r.client.Indices.Create.WithBody(bytes.NewReader(body)),
		r.client.Indices.Create.WithContext(ctx),
	)
	if err != nil {
		return err
	}
	defer createRes.Body.Close()

	if createRes.StatusCode >= 300 {
		b, _ := io.ReadAll(createRes.Body)
		return appErrors.WrapperElasticsearchError(createRes.StatusCode, "create index failed", string(b))
	}

	return nil
}

// DeleteIndex deletes an index
func (r *Repository[T]) DeleteIndex(ctx context.Context, index string) error {
	res, err := r.client.Indices.Delete([]string{index}, r.client.Indices.Delete.WithContext(ctx))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return appErrors.ErrIndexNotFound
	}
	if res.StatusCode >= 300 {
		b, _ := io.ReadAll(res.Body)
		return appErrors.WrapperElasticsearchError(res.StatusCode, "delete index failed", string(b))
	}

	return nil
}

// Refresh refreshes an index
func (r *Repository[T]) Refresh(ctx context.Context, index string) error {
	res, err := r.client.Indices.Refresh(
		r.client.Indices.Refresh.WithIndex([]string{index}...),
		r.client.Indices.Refresh.WithContext(ctx),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode >= 300 {
		b, _ := io.ReadAll(res.Body)
		return appErrors.WrapperElasticsearchError(res.StatusCode, "refresh failed", string(b))
	}

	return nil
}

// Index inserts or updates a document
func (r *Repository[T]) Index(ctx context.Context, index, id string, doc T) error {
	b, err := json.Marshal(doc)
	if err != nil {
		return err
	}

	res, err := r.client.Index(
		index,
		bytes.NewReader(b),
		r.client.Index.WithDocumentID(id),
		r.client.Index.WithContext(ctx),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode >= 300 {
		body, _ := io.ReadAll(res.Body)
		return appErrors.WrapperElasticsearchError(res.StatusCode, "index failed", string(body))
	}

	return nil
}

// BulkIndex indexes multiple documents at once
func (r *Repository[T]) BulkIndex(ctx context.Context, index string, docs []T, idSelector func(T) string) error {
	if len(docs) == 0 {
		return nil
	}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)

	for _, d := range docs {
		meta := map[string]map[string]string{"index": {"_index": index}}
		if idSelector != nil {
			meta["index"]["_id"] = idSelector(d)
		}
		if err := enc.Encode(meta); err != nil {
			return err
		}
		if err := enc.Encode(d); err != nil {
			return err
		}
	}

	res, err := r.client.Bulk(bytes.NewReader(buf.Bytes()), r.client.Bulk.WithIndex(index), r.client.Bulk.WithContext(ctx))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode >= 300 {
		b, _ := io.ReadAll(res.Body)
		return appErrors.WrapperElasticsearchError(res.StatusCode, "bulk failed", string(b))
	}

	var br types.BulkResponse
	if err := json.NewDecoder(res.Body).Decode(&br); err == nil {
		if br.Errors {
			return appErrors.ErrBulkFailed
		}
	}

	return nil
}

// Search executes a query and returns results
func (r *Repository[T]) Search(ctx context.Context, index string, query map[string]any, from, size int, sort []map[string]any) ([]T, int64, error) {
	reqMap := map[string]any{}
	if query != nil {
		reqMap["query"] = query
	} else {
		reqMap["query"] = map[string]any{"match_all": map[string]any{}}
	}
	if size > 0 {
		reqMap["size"] = size
	}
	if from > 0 {
		reqMap["from"] = from
	}
	if len(sort) > 0 {
		reqMap["sort"] = sort
	}

	body, err := json.Marshal(reqMap)
	if err != nil {
		return nil, 0, err
	}

	res, err := r.client.Search(
		r.client.Search.WithContext(ctx),
		r.client.Search.WithIndex(index),
		r.client.Search.WithBody(bytes.NewReader(body)),
	)
	if err != nil {
		return nil, 0, err
	}
	defer res.Body.Close()

	if res.StatusCode >= 300 {
		b, _ := io.ReadAll(res.Body)
		return nil, 0, appErrors.WrapperElasticsearchError(res.StatusCode, "search failed", string(b))
	}

	var raw struct {
		Hits struct {
			Total struct {
				Value int64 `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source T `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&raw); err != nil {
		return nil, 0, err
	}

	out := make([]T, 0, len(raw.Hits.Hits))
	for _, h := range raw.Hits.Hits {
		out = append(out, h.Source)
	}

	return out, raw.Hits.Total.Value, nil
}
