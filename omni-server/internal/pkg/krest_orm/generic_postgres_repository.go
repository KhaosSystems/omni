package krest_orm

/*
* Contains a simple impmentation of a krest repository. Meant to be used as reference for implementing a repository,
* a simple repo for testing, or as a base.
 */

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
)

// Implement the krest.Repository[T] interface.
type GenericPostgresRepository[T any] struct {
	db    *sql.DB
	table string
}

func NewGenericPostgresRepository[T any](db *sql.DB) *GenericPostgresRepository[T] {
	table := TableName[T]()
	schema := Schema[T]()

	// Check if table exists, and matches schema.
	// TODO: Throw error if schema does not match.
	// TODO: Add a migration system.
	sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s)", table, schema)
	log.Printf("creating table: %s", sql)
	_, err := db.Exec(sql)
	if err != nil {
		log.Fatalf("failed to create table %s: %v", table, err)
	}

	return &GenericPostgresRepository[T]{
		db:    db,
		table: table,
	}
}

func (r *GenericPostgresRepository[T]) Get(ctx context.Context, id uuid.UUID, query krest.ResourceQuery) (T, error) {
	// Initialize a new value of T
	var t T
	tValue := reflect.ValueOf(&t).Elem()
	tType := reflect.TypeOf(tValue)
	if tType.Kind() != reflect.Struct {
		return *new(T), fmt.Errorf("type %T is not a struct", t)
	}

	// Fields to get from the database.
	var fieldsToGet []string

	// We always need the non-expandable fields.
	nonExpandableFields, err := krest.ReflectNonExpandableFields[T]()
	if err != nil {
		return *new(T), fmt.Errorf("failed to get non-expandable fields: %v", err)
	}
	for key := range nonExpandableFields {
		fieldsToGet = append(fieldsToGet, key)
	}

	// Get the expandable fields.
	expandableFields, err := krest.ReflectExpandableFields[T]()
	if err != nil {
		return *new(T), fmt.Errorf("failed to get expandable fields: %v", err)
	}
	for key := range expandableFields {
		// Add if requested.
		if slices.Contains(query.Expand, key) {
			fieldsToGet = append(fieldsToGet, key)
		}
	}

	// Ensure that at least one field is selected
	if len(fieldsToGet) == 0 {
		return *new(T), fmt.Errorf("no fields selected for query.. not event the uuid.. for some reason")
	}

	// Get the fields from the database.
	queryFields := strings.Join(fieldsToGet, ", ")
	sql := fmt.Sprintf("SELECT %s FROM tasks WHERE uuid = $1", queryFields)

	// Execute the query.
	row := r.db.QueryRow(sql, id)

	// Scan the row into the task struct.
	var task T
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
			return *new(T), fmt.Errorf("field %s not found in task struct", field)
		}
	}

	log.Printf("dest fields: %v", destFields)

	// Scan the row into the task struct.
	if err := row.Scan(destFields...); err != nil {
		return *new(T), fmt.Errorf("failed to scan task: %v", err)
	}

	// The task should now be populated with the selected fields.
	return task, nil
}

func (r *GenericPostgresRepository[T]) List(ctx context.Context, query krest.CollectionQuery) ([]T, error) {
	// Initialize a new value of T
	var t T
	tValue := reflect.ValueOf(&t).Elem()
	tType := reflect.TypeOf(tValue)
	if tType.Kind() != reflect.Struct {
		return []T{}, fmt.Errorf("type %T is not a struct", t)
	}

	// Fields to get from the database.
	var fieldsToGet []string

	// We always need the non-expandable fields.
	nonExpandableFields, err := krest.ReflectNonExpandableFields[T]()
	if err != nil {
		return []T{}, fmt.Errorf("failed to get non-expandable fields: %v", err)
	}
	for key := range nonExpandableFields {
		fieldsToGet = append(fieldsToGet, key)
	}

	// Get the expandable fields.
	expandableFields, err := krest.ReflectExpandableFields[T]()
	if err != nil {
		return []T{}, fmt.Errorf("failed to get expandable fields: %v", err)
	}
	for key := range expandableFields {
		// Add if requested.
		if slices.Contains(query.Expand, key) {
			fieldsToGet = append(fieldsToGet, key)
		}
	}

	// Ensure that at least one field is selected
	if len(fieldsToGet) == 0 {
		return []T{}, fmt.Errorf("no fields selected for query.. not event the uuid.. for some reason")
	}

	// Get the fields from the database.
	queryFields := strings.Join(fieldsToGet, ", ")
	sql := fmt.Sprintf("SELECT %s, COUNT(*) OVER() AS total FROM %s", queryFields, r.table)
	args := []interface{}{}
	argIdx := 1

	// Add the limit and offset to the query.
	if query.Limit > 0 {
		sql += fmt.Sprintf(" LIMIT $%d", argIdx)
		args = append(args, query.Limit)
		argIdx++
	}

	if query.Offset > 0 {
		sql += fmt.Sprintf(" OFFSET $%d", argIdx)
		args = append(args, query.Offset)
	}

	// Query the database for all projects.
	log.Printf("query: %s, args: %v", sql, args)
	rows, err := r.db.Query(sql, args...)
	if err != nil {
		return []T{}, err
	}
	defer rows.Close()

	// Scan the row into the task struct.
	var task T
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
			return []T{}, fmt.Errorf("field %s not found in task struct", field)
		}
	}
	log.Printf("dest fields: %v", destFields)

	// Iterate over the rows and create a slice of projects.
	tasks := []T{}
	var total int
	destFields = append(destFields, &total) // add the &total
	for rows.Next() {
		var resource T
		err := rows.Scan(destFields...)
		if err != nil {
			return []T{}, err
		}
		tasks = append(tasks, resource)
	}

	// The task should now be populated with the selected fields.
	return tasks, nil
}

func (r *GenericPostgresRepository[T]) Create(ctx context.Context, user T) (T, error) {
	return *new(T), errors.ErrUnsupported
}

func (r *GenericPostgresRepository[T]) Update(ctx context.Context, id uuid.UUID, user T) (T, error) {
	return *new(T), errors.ErrUnsupported
}

func (r *GenericPostgresRepository[T]) Delete(ctx context.Context, id uuid.UUID) error {
	return errors.ErrUnsupported
}
