package task

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"reflect"
	"slices"
	"strings"

	"github.com/google/uuid"
	"github.com/khaossystems/omni-server/internal/pkg/krest"
	"github.com/khaossystems/omni-server/pkg/models"
)

type TaskRepository = krest.Repository[models.Task]

// Implement the TaskRepository interface.
type PostgresTaskRepository struct {
	db *sql.DB
}

func NewPostgresTaskRepository(db *sql.DB) *PostgresTaskRepository {
	return &PostgresTaskRepository{db: db}
}

func (r *PostgresTaskRepository) Get(ctx context.Context, id uuid.UUID, query krest.ResourceQuery) (models.Task, error) {
	// Fields to get from the database.
	var fieldsToGet []string

	// We always need the non-expandable fields.
	nonExpandableFields, err := krest.ReflectNonExpandableFields[models.Task]()
	if err != nil {
		return models.Task{}, fmt.Errorf("failed to get non-expandable fields: %v", err)
	}
	for key := range nonExpandableFields {
		fieldsToGet = append(fieldsToGet, key)
	}

	// Get the expandable fields.
	expandableFields, err := krest.ReflectExpandableFields[models.Task]()
	if err != nil {
		return models.Task{}, fmt.Errorf("failed to get expandable fields: %v", err)
	}
	for key := range expandableFields {
		// Add if requested.
		if slices.Contains(query.Expand, key) {
			fieldsToGet = append(fieldsToGet, key)
		}
	}

	// Ensure that at least one field is selected
	if len(fieldsToGet) == 0 {
		return models.Task{}, fmt.Errorf("no fields selected for query.. not event the uuid.. for some reason")
	}

	// Get the fields from the database.
	queryFields := strings.Join(fieldsToGet, ", ")
	sql := fmt.Sprintf("SELECT %s FROM tasks WHERE uuid = $1", queryFields)

	// Execute the query.
	row := r.db.QueryRow(sql, id)

	// Scan the row into the task struct.
	var task models.Task
	destFields := make([]interface{}, len(fieldsToGet))

	taskValue := reflect.ValueOf(&task).Elem() // Get the value of the task struct.

	log.Printf("getting fields: %v", fieldsToGet)

	for i, field := range fieldsToGet {
		// Match the field with the struct field by name.
		structField := taskValue.FieldByNameFunc(func(s string) bool {
			// Match struct field name with the database field (case-sensitive check).
			return strings.EqualFold(s, field)
		})

		if structField.IsValid() {
			destFields[i] = structField.Addr().Interface() // Get address to assign value via Scan.
		} else {
			return models.Task{}, fmt.Errorf("field %s not found in task struct", field)
		}
	}

	log.Printf("dest fields: %v", destFields)

	// Scan the row into the task struct.
	if err := row.Scan(destFields...); err != nil {
		return models.Task{}, fmt.Errorf("failed to scan task: %v", err)
	}

	// The task should now be populated with the selected fields.
	return task, nil
}

func (r *PostgresTaskRepository) List(ctx context.Context, query krest.CollectionQuery) ([]models.Task, error) {
	return []models.Task{}, nil
}

func (r *PostgresTaskRepository) Create(ctx context.Context, user models.Task) (models.Task, error) {
	return models.Task{}, errors.ErrUnsupported
}

func (r *PostgresTaskRepository) Update(ctx context.Context, id uuid.UUID, user models.Task) (models.Task, error) {
	return models.Task{}, errors.ErrUnsupported
}

func (r *PostgresTaskRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return errors.ErrUnsupported
}
