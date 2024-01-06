package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// Cols holds the schema definition for the Cols entity.
type Cols struct {
	ent.Schema
}

func (Cols) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "columns",
		},
	}
}

// Fields of the Cols.
func (Cols) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Default(0),
		field.String("table_schema").Default("").MaxLen(64),
		field.String("table_name").Default("").MaxLen(64),
		field.String("column_name").Default("").MaxLen(64),
		// field.Text("column_default").Default(""),
		// field.Text("column_comment").Default(""),
	}
}

// Edges of the Cols.
func (Cols) Edges() []ent.Edge {
	return nil
}
