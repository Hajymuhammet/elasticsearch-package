package usecase

import (
	"context"
	"errors"

	"github.com/salamsites/elasticsearch-package/pkg/mapping"
	"github.com/salamsites/elasticsearch-package/repository"
)

// Service[T] is a generic Elasticsearch service wrapping repository operations.
type Service[T any] struct {
	repo repository.IndexRepository[T]
}

// NewService creates a new Service instance
func NewService[T any](repo repository.IndexRepository[T]) *Service[T] {
	return &Service[T]{repo: repo}
}

// EnsureIndex ensures an index exists, creating it with the provided mapping if missing
func (s *Service[T]) EnsureIndex(ctx context.Context, index string, mappingBody map[string]any) error {
	if index == "" {
		return errors.New("index is empty")
	}
	return s.repo.EnsureIndex(ctx, index, mappingBody)
}

// DeleteIndex deletes an index
func (s *Service[T]) DeleteIndex(ctx context.Context, index string) error {
	if index == "" {
		return errors.New("index required")
	}
	return s.repo.DeleteIndex(ctx, index)
}

// Refresh refreshes an index
func (s *Service[T]) Refresh(ctx context.Context, index string) error {
	if index == "" {
		return errors.New("index required")
	}
	return s.repo.Refresh(ctx, index)
}

// Index indexes a single document
func (s *Service[T]) Index(ctx context.Context, index, id string, doc T) error {
	if index == "" || id == "" {
		return errors.New("index and id are required")
	}
	return s.repo.Index(ctx, index, id, doc)
}

// BulkIndex indexes multiple documents at once
func (s *Service[T]) BulkIndex(ctx context.Context, index string, docs []T, idSelector func(T) string) error {
	if index == "" {
		return errors.New("index required")
	}
	if len(docs) == 0 {
		return nil
	}
	if idSelector == nil {
		return errors.New("idSelector required")
	}
	return s.repo.BulkIndex(ctx, index, docs, idSelector)
}

// Search executes a search query
func (s *Service[T]) Search(ctx context.Context, index string, query map[string]any, from, size int, sort []map[string]any) ([]T, int64, error) {
	if index == "" {
		return nil, 0, errors.New("index required")
	}
	return s.repo.Search(ctx, index, query, from, size, sort)
}

func (s *Service[T]) DeleteDocument(ctx context.Context, index, id string) error {
	if index == "" || id == "" {
		return errors.New("index and id required")
	}
	return s.repo.DeleteDocument(ctx, index, id)
}

// GetByID fetches a document by ID
func (s *Service[T]) GetByID(ctx context.Context, index, id string) (*T, error) {
	if index == "" || id == "" {
		return nil, errors.New("index and id required")
	}
	return s.repo.GetByID(ctx, index, id)
}

// BuildMapping generates an Elasticsearch mapping from a struct using reflection and tags
func (s *Service[T]) BuildMapping(sample T) map[string]any {
	return mapping.BuildMappingFromStruct(sample)
}
