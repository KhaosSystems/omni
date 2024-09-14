package sql

import (
	"fmt"
	"strings"
)

/*
* TableSchema is a helper struct representing a SQL table schema.
 */
type TableSchema struct {
	Name    string
	Columns []ColumnSchema
}

func (s TableSchema) ColumnDefinitions() []string {
	definitions := []string{}
	for _, column := range s.Columns {
		definitions = append(definitions, column.String())
	}
	return definitions
}

func (s TableSchema) CreateTableQuery() string {
	columns := strings.Join(s.ColumnDefinitions(), ", ")
	return fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s)", s.Name, columns)
}

/*
* TableSchemaBuilder is a helper struct for building, and validating, SQL schemas.
 */
type TableSchemaBuilder struct {
	name    string
	columns []ColumnSchema
}

func NewTableSchemaBuilder() *TableSchemaBuilder {
	return &TableSchemaBuilder{}
}

func (b *TableSchemaBuilder) Name(name string) *TableSchemaBuilder {
	b.name = name
	return b
}

func (b *TableSchemaBuilder) Column(schema ColumnSchema) *TableSchemaBuilder {
	b.columns = append(b.columns, schema)
	return b
}

func (b *TableSchemaBuilder) Build() (TableSchema, error) {
	if b.name == "" {
		return TableSchema{}, fmt.Errorf("table name is required")
	}

	if len(b.columns) == 0 {
		return TableSchema{}, fmt.Errorf("at least one column is required (table: %s)", b.name)
	}

	return TableSchema{
		Name:    b.name,
		Columns: b.columns,
	}, nil
}
