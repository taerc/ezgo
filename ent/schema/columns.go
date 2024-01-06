package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// information_schema
// Columns holds the schema definition for the Columns entity.
type Columns struct {
	ent.Schema
}

// Annotations of the Admin.
func (Columns) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "COLUMNS",
		},
	}
}

// Fields of the Columns.
func (Columns) Fields() []ent.Field {
	return []ent.Field{
		field.String("TABLE_SCHEMA").Default("").MaxLen(64),
		field.String("TABLE_NAME").Default("").MaxLen(64),
		field.String("COLUMN_NAME").Default("").MaxLen(64),
		field.Text("COLUMN_DEFAULT").Default(""),
		field.Text("COLUMN_COMMENT").Default(""),
	}
}

// Edges of the Columns.
func (Columns) Edges() []ent.Edge {
	return nil
}
