package repository

import "context"

type IndexRepository[T any] interface {
	EnsureIndex(ctx context.Context, index string, mapping map[string]any) error
	DeleteIndex(ctx context.Context, index string) error
	Refresh(ctx context.Context, index string) error
	Index(ctx context.Context, index, id string, doc T) error
	BulkIndex(ctx context.Context, index string, docs []T, idSelector func(T) string) error
	Search(ctx context.Context, index string, query map[string]any, from, size int, sort []map[string]any) ([]T, int64, error)
	DeleteDocument(ctx context.Context, index, id string) error
	GetByID(ctx context.Context, index, id string) (*T, error)
}
