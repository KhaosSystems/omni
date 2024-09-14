package sql

/*
* TableSchema is a helper struct representing a SQL table schema.
 */
type TableSchema struct {
	Name    string
	Columns []ColumnSchema
}

func (s TableSchema) String() string {
	return ""
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

func (b *TableSchemaBuilder) Column(schema ColumnSchema) error {
	return nil
}

func (b *TableSchemaBuilder) Build() (TableSchema, error) {
	return TableSchema{
		Name:    b.name,
		Columns: b.columns,
	}, nil
}
