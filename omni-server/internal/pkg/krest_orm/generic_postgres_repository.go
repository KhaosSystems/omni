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
	"github.com/jmoiron/sqlx"
	"github.com/khaossystems/omni-server/internal/pkg/krest"
)

// Implement the krest.Repository[T] interface.
type GenericPostgresRepository[T any] struct {
	db    *sqlx.DB
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
		db:    sqlx.NewDb(db, "postgres"),
		table: table,
	}
}

/*
* Returns the fields for a given struct type and query.
 */
func (r *GenericPostgresRepository[T]) FieldsToGet(expand []string) ([]reflect.StructField, error) {
	var t T
	tValue := reflect.ValueOf(&t).Elem()
	tType := reflect.TypeOf(tValue)

	// Fields to get from the database.
	fields := []reflect.StructField{}

	// Make sure T is a struct.
	if tType.Kind() != reflect.Struct {
		return []reflect.StructField{}, fmt.Errorf("type %T is not a struct", t)
	}

	// We always need the non-expandable fields.
	nonExpandableFields, err := krest.ReflectNonExpandableFields[T]()
	if err != nil {
		return []reflect.StructField{}, fmt.Errorf("failed to get non-expandable fields: %v", err)
	}
	for _, field := range nonExpandableFields {
		fields = append(fields, field)
	}

	// Get the expandable fields.
	expandableFields, err := krest.ReflectExpandableFields[T]()
	if err != nil {
		return []reflect.StructField{}, fmt.Errorf("failed to get expandable fields: %v", err)
	}
	for _, field := range expandableFields {
		// Add if requested.
		fieldColumnName := ColumnName(field.Name)
		if slices.Contains(expand, fieldColumnName) {
			fields = append(fields, field)
		}
	}

	// Ensure that at least one field is selected
	if len(fields) == 0 {
		return []reflect.StructField{}, fmt.Errorf("no fields selected for query.. not event the uuid.. for some reason")
	}

	return fields, nil
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
	fieldsToGet, err := r.FieldsToGet(query.Expand)
	if err != nil {
		return *new(T), fmt.Errorf("failed to get fields to get: %v", err)
	}

	columnNamesToGet := []string{}
	for _, field := range fieldsToGet {
		columnNamesToGet = append(columnNamesToGet, ColumnName(field.Name))
	}

	// Get the fields from the database.
	queryFields := strings.Join(columnNamesToGet, ", ")
	sql := fmt.Sprintf("SELECT %s FROM tasks WHERE uuid = $1", queryFields)

	// Execute the query.
	resource := *new(T)
	err = r.db.Select(resource, sql, id)
	if err != nil {
		return *new(T), fmt.Errorf("failed to query database: %v", err)
	}

	return resource, nil
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
	fieldsToGet, err := r.FieldsToGet(query.Expand)
	if err != nil {
		return []T{}, fmt.Errorf("failed to get fields to get: %v", err)
	}

	// Get the db names of the fields to get.
	columnNamesToGet := []string{}
	for _, field := range fieldsToGet {
		columnNamesToGet = append(columnNamesToGet, ColumnName(field.Name))
	}

	// Get the fields from the database.
	queryFields := strings.Join(columnNamesToGet, ", ")
	sql := fmt.Sprintf("SELECT %s FROM %s", queryFields, r.table)
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
	log.Printf("query: %s, args: %v", sql, args)

	// Query the database for all projects.
	resources := []T{}
	err = r.db.Select(&resources, sql, args...)
	if err != nil {
		return []T{}, fmt.Errorf("failed to query database: %v", err)
	}

	// Get the total count of resources.
	var total int
	totalQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s", r.table)
	err = r.db.Get(&total, totalQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get total users: %v", err)
	}

	// The task should now be populated with the selected fields.
	return resources, nil
}

func (r *GenericPostgresRepository[T]) Create(ctx context.Context, resource T) (T, error) {
	// Use reflection to get the value and type of the resource
	resourceValue := reflect.ValueOf(&resource).Elem()
	resourceType := resourceValue.Type()

	if resourceType.Kind() != reflect.Struct {
		return *new(T), fmt.Errorf("type %T is not a struct", resource)
	}

	// Prepare slices for column names and placeholders
	columnNames := []string{}
	placeholders := []string{}
	values := []interface{}{}
	argIdx := 1

	// Iterate over the struct fields
	for i := 0; i < resourceValue.NumField(); i++ {
		field := resourceType.Field(i)
		fieldValue := resourceValue.Field(i)

		// Skip unexported fields
		if !fieldValue.CanInterface() {
			continue
		}

		// Use the field name as the column name or map to a proper DB column name using a helper
		columnName := ColumnName(field.Name)

		// Append column names and placeholders
		columnNames = append(columnNames, columnName)
		placeholders = append(placeholders, fmt.Sprintf("$%d", argIdx))
		values = append(values, fieldValue.Interface())
		argIdx++
	}

	// Generate the SQL query
	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s) RETURNING *",
		r.table,
		strings.Join(columnNames, ", "),
		strings.Join(placeholders, ", "),
	)

	// Execute the query and return the created resource
	var createdResource T
	err := r.db.GetContext(ctx, &createdResource, query, values...)
	if err != nil {
		return *new(T), fmt.Errorf("failed to insert resource into database: %v", err)
	}

	return createdResource, nil
}

func (r *GenericPostgresRepository[T]) Update(ctx context.Context, id uuid.UUID, user T) (T, error) {
	return *new(T), errors.ErrUnsupported
}

func (r *GenericPostgresRepository[T]) Delete(ctx context.Context, id uuid.UUID) error {
	return errors.ErrUnsupported
}
