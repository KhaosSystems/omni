package krest_orm

import (
	"context"

	"github.com/google/uuid"
	"github.com/khaossystems/omni-server/internal/pkg/krest"
)

// Implements krest.Service[T]
type GenericService[T any] struct {
	repository krest.Repository[T]
}

func NewGenericService[T any](repository krest.Repository[T]) *GenericService[T] {
	return &GenericService[T]{repository: repository}
}

func (s *GenericService[T]) Get(ctx context.Context, id uuid.UUID, query krest.ResourceQuery) (T, error) {
	return s.repository.Get(ctx, id, query)
}

func (s *GenericService[T]) List(ctx context.Context, query krest.CollectionQuery) ([]T, error) {
	return s.repository.List(ctx, query)
}

func (s *GenericService[T]) Create(ctx context.Context, user T) (T, error) {
	return s.repository.Create(ctx, user)
}

func (s *GenericService[T]) Update(ctx context.Context, id uuid.UUID, user T) (T, error) {
	return s.repository.Update(ctx, id, user)
}

func (s *GenericService[T]) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repository.Delete(ctx, id)
}
