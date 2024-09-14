package sql

import (
	"fmt"
	"strings"
)

/*
* ColumnSchema is a helper struct representing a SQL column schema.
 */
type ColumnSchema struct {
	Name        string
	Type        string
	Constraints []string
}

func (c ColumnSchema) String() string {
	constraints := strings.Join(c.Constraints, " ")
	return fmt.Sprintf("%s %s %s", c.Name, c.Type, constraints)
}

/*
* ColumnSchemaBuilder is a fluent builder for building, and validating, SQL column schemas.
 */
type ColumnSchemaBuilder struct {
	name        string
	sqlType     string
	constraints []string
}

func NewColumnSchemaBuilder() *ColumnSchemaBuilder {
	return &ColumnSchemaBuilder{}
}

func (b *ColumnSchemaBuilder) Name(name string) *ColumnSchemaBuilder {
	b.name = name
	return b
}

func (b *ColumnSchemaBuilder) Type(sqlType string) *ColumnSchemaBuilder {
	b.sqlType = sqlType
	return b
}

func (b *ColumnSchemaBuilder) AddConstraint(constraint string) *ColumnSchemaBuilder {
	b.constraints = append(b.constraints, constraint)
	return b
}

func (b *ColumnSchemaBuilder) AddConstraints(constraints ...string) *ColumnSchemaBuilder {
	b.constraints = append(b.constraints, constraints...)
	return b
}

func (b *ColumnSchemaBuilder) Build() (ColumnSchema, error) {
	// TODO: Validate the schema.
	return ColumnSchema{
		Name:        b.name,
		Type:        b.sqlType,
		Constraints: b.constraints,
	}, nil
}
