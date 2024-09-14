package sql

import (
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/gertd/go-pluralize"
	"github.com/google/uuid"
	"github.com/khaossystems/omni-server/internal/pkg/krest"
)

/*
* Returns all krest_orm tags for a given field.
 */
func GetKrestTags(field reflect.StructField) []string {
	tags := field.Tag.Get("krest_orm")
	//fmt.Printf("%s: %s", field.Name, tags)
	return strings.Split(tags, ",")
}

/*
* Returns the column name for a field name.
* By default, the column name is the snake_case version of the field name.
 */
func ColumnName(fieldName string) string {
	return krest.ToSnakeCase(fieldName)
}

/*
* Returns the SQL table name for a given struct type.
* By default, the table name is an all lowercase, snake_case version of the type name.
* TODO: Pluralize the table name.
 */
func TableName[T any]() string {
	var t T
	snake := krest.ToSnakeCase(reflect.TypeOf(t).Name())

	// Pluralize the table name.
	pluralize := pluralize.NewClient()
	plural := pluralize.Plural(snake)

	return plural
}

/*
* Helper function for converting go types to SQL types.
* TODO: Allow for custom types using the krest_orm tag (type:varchar(36)).
 */
func GoTypeToSQLType(goType reflect.Type) (string, error) {
	switch goType.Kind() {
	case reflect.String:
		return "TEXT", nil
	case reflect.Int:
		return "INTEGER", nil
	case reflect.Int64:
		return "BIGINT", nil
	case reflect.Bool:
		return "BOOLEAN", nil
	case reflect.Float64:
		return "DOUBLE PRECISION", nil
	case reflect.Array:
		// uuid.UUID
		if goType == reflect.TypeOf(uuid.UUID{}) {
			return "UUID", nil
		}
	}

	return "", fmt.Errorf("unsupported type: %s, kind: %s", goType, goType.Kind())
}

/*
* Helper function for converting a struct field to a SQL column schema.
 */
func ColumnSchemaFromField(field reflect.StructField) (ColumnSchema, error) {
	builder := NewColumnSchemaBuilder()

	// Name.
	builder.Name(ColumnName(field.Name))

	// Find the SQL type of the field.
	// TODO: Allow for custom types using the krest_orm tag (type:varchar(36)).
	sqlType, err := GoTypeToSQLType(field.Type)
	if err != nil {
		return ColumnSchema{}, err
	}

	builder.Type(sqlType)

	// Constraints.
	tags := GetKrestTags(field)
	if slices.Contains(tags, "pk") {
		builder.AddConstraint("PRIMARY KEY")
	}

	if slices.Contains(tags, "fk") {
		builder.AddConstraint("FOREIGN KEY")
	}

	// Build.
	return builder.Build()
}

/*
* Helper function for converting a struct to a SQL table schema
 */
func TableSchemaFromStruct[T any]() (TableSchema, error) {
	builder := NewTableSchemaBuilder()

	// Get the struct type.
	var t T
	tType := reflect.TypeOf(t)

	// Make sure the type is a struct.
	if tType.Kind() != reflect.Struct {
		return TableSchema{}, fmt.Errorf("type %T is not a struct", t)
	}

	// Name.
	builder.Name(TableName[T]())

	// Columns.
	for i := 0; i < tType.NumField(); i++ {
		field := tType.Field(i)

		// Skip if tagged as ignored.
		if slices.Contains(GetKrestTags(field), "ignore") {
			continue
		}

		// Convert the field to a column schema.
		colSchema, err := ColumnSchemaFromField(field)
		if err != nil {
			return TableSchema{}, err
		}

		builder.Column(colSchema)
	}

	return builder.Build()
}

// EnsurePrimaryKey sets the primary key field to a new UUID if it's not already set.
func EnsurePrimaryKey[T any](resource T) (T, error) {
	// Reflect on the resource to access its fields.
	resourceValue := reflect.ValueOf(&resource).Elem()

	// Get the primary key field.
	primaryKeyField, err := krest.ReflectPrimaryKeyField[T]()
	if err != nil {
		return resource, fmt.Errorf("failed to get primary key field: %v", err)
	}

	// Ensure the primary key field is of the expected type.
	if primaryKeyField.Type != reflect.TypeOf(uuid.UUID{}) {
		return resource, fmt.Errorf("unsupported primary key type: %s", primaryKeyField.Type)
	}

	// Get the current value of the primary key field.
	currentValue := resourceValue.FieldByName(primaryKeyField.Name).Interface()

	// If the current value is a zero value (uuid.Nil), set it to a new UUID.
	if currentValue == uuid.Nil {
		resourceValue.FieldByName(primaryKeyField.Name).Set(reflect.ValueOf(uuid.New()))
	}

	return resource, nil
}
