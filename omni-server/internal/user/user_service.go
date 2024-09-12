package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/khaossystems/omni-server/internal/pkg/krest"
	"github.com/khaossystems/omni-server/pkg/models"
)

// Implements krest.Service[models.User]
type UserService struct {
	repository UserRepository
}

func NewUserService(repository UserRepository) *UserService {
	return &UserService{repository: repository}
}

func (s *UserService) Get(ctx context.Context, id uuid.UUID, query krest.ResourceQuery) (models.User, error) {
	return s.repository.Get(ctx, id, query)
}

func (s *UserService) List(ctx context.Context, query krest.CollectionQuery) ([]models.User, error) {
	return s.repository.List(ctx, query)
}

func (s *UserService) Create(ctx context.Context, user models.User) (models.User, error) {
	return s.repository.Create(ctx, user)
}

func (s *UserService) Update(ctx context.Context, id uuid.UUID, user models.User) (models.User, error) {
	return s.repository.Update(ctx, id, user)
}

func (s *UserService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repository.Delete(ctx, id)
}
