package krest_orm

import (
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"strings"

	pluralize "github.com/gertd/go-pluralize"
	"github.com/google/uuid"
)

/*
* Helper function for converting a string to snake_case.
 */
func ToSnakeCase(str string) string {
	matchFirstCap := regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap := regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
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
	return ToSnakeCase(fieldName)
}

/*
* Helper function for converting a field to a SQL column.
 */
func FieldToColumn(field reflect.StructField) string {
	tags := GetKrestTags(field)
	if slices.Contains(tags, "ignore") {
		return ""
	}

	columnName := ColumnName(field.Name)
	typeName, err := GoTypeToSQLType(field.Type)
	if err != nil {
		panic(err)
	}

	if slices.Contains(tags, "pk") {
		return fmt.Sprintf("%s %s PRIMARY KEY", columnName, typeName)
	}

	return fmt.Sprintf("%s %s", columnName, typeName)
}

/*
* Returns the SQL table name for a given struct type.
* By default, the table name is an all lowercase, snake_case version of the type name.
* TODO: Pluralize the table name.
 */
func TableName[T any]() string {
	var t T
	snake := ToSnakeCase(reflect.TypeOf(t).Name())

	// Pluralize the table name.
	pluralize := pluralize.NewClient()
	plural := pluralize.Plural(snake)

	return plural
}

/*
* Generates a SQL schema given a struct type.
* By default, the column name are all lowercase, snake_case versions of the struct field names.
 */
func Schema[T any]() string {
	var t T
	tType := reflect.TypeOf(t)

	if tType.Kind() != reflect.Struct {
		panic("type T is not a struct")
	}

	var schema string
	for i := 0; i < tType.NumField(); i++ {
		field := tType.Field(i)
		colSchema := FieldToColumn(field)
		if colSchema != "" {
			schema += colSchema + ", "
		}
	}

	// Remove the trailing comma and space.
	if len(schema) > 1 {
		schema = schema[:len(schema)-2]
	}

	return schema
}
