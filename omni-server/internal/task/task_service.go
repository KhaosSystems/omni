package task

import (
	"context"

	"github.com/google/uuid"
	"github.com/khaossystems/omni-server/internal/pkg/krest"
	"github.com/khaossystems/omni-server/pkg/models"
)

// Implements krest.Service[models.Task]
type TaskService struct {
	repository TaskRepository
}

func NewTaskService(repository TaskRepository) *TaskService {
	return &TaskService{repository: repository}
}

func (s *TaskService) Get(ctx context.Context, id uuid.UUID, query krest.ResourceQuery) (models.Task, error) {
	return s.repository.Get(ctx, id, query)
}

func (s *TaskService) List(ctx context.Context, query krest.CollectionQuery) ([]models.Task, error) {
	return s.repository.List(ctx, query)
}

func (s *TaskService) Create(ctx context.Context, user models.Task) (models.Task, error) {
	return s.repository.Create(ctx, user)
}

func (s *TaskService) Update(ctx context.Context, id uuid.UUID, user models.Task) (models.Task, error) {
	return s.repository.Update(ctx, id, user)
}

func (s *TaskService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repository.Delete(ctx, id)
}
