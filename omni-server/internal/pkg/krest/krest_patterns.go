package krest

import (
	"context"

	"github.com/google/uuid"
)

/*
* Repository interface for use with the repository pattern. Defines basic CRUD operations.
* Note: We recommend using this pattens for consitancy, but it's by no means required.
 */
type Repository[T any] interface {
	Get(ctx context.Context, id uuid.UUID, query ResourceQuery) (T, error)
	List(ctx context.Context, query CollectionQuery) ([]T, error)
	Create(ctx context.Context, resource T) (T, error)
	Update(ctx context.Context, id uuid.UUID, resource T) (T, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

/*
* Service interface for use with the service pattern. Defines basic CRUD operations.
* Note: We recommend using this pattens for consitancy, but it's by no means required.
 */
type Service[T any] interface {
	Get(ctx context.Context, id uuid.UUID, query ResourceQuery) (T, error)
	List(ctx context.Context, query CollectionQuery) ([]T, error)
	Create(ctx context.Context, resource T) (T, error)
	Update(ctx context.Context, id uuid.UUID, resource T) (T, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
