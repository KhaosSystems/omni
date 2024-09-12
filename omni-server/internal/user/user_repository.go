package user

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/khaossystems/omni-server/internal/pkg/krest"
	"github.com/khaossystems/omni-server/pkg/models"
)

// constant mock user array
var users = []models.User{
	{UUID: uuid.New(), Name: "Alice"},
	{UUID: uuid.New(), Name: "Bob"},
}

type UserRepository = krest.Repository[models.User]

// Implement the UserRepository interface.
type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) Get(ctx context.Context, id uuid.UUID, query krest.ResourceQuery) (models.User, error) {
	// Find the user with the given ID.
	for _, user := range users {
		if user.UUID == id {
			return user, nil
		}
	}

	// Return an error if the user was not found.
	return models.User{}, errors.New("user not found")
}

func (r *PostgresUserRepository) List(ctx context.Context, query krest.CollectionQuery) ([]models.User, error) {
	return users, nil
}

func (r *PostgresUserRepository) Create(ctx context.Context, user models.User) (models.User, error) {
	return models.User{}, errors.ErrUnsupported
}

func (r *PostgresUserRepository) Update(ctx context.Context, id uuid.UUID, user models.User) (models.User, error) {
	return models.User{}, errors.ErrUnsupported
}

func (r *PostgresUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return errors.ErrUnsupported
}
