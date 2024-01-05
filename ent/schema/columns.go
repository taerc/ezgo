package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema"
)

// Columns holds the schema definition for the Columns entity.
type Columns struct {
	ent.Schema
}

// Fields of the Columns.
func (Columns) Fields() []ent.Field {
	return []schema.Annotation{ensql.Annotation{Table: "COLUMNS"}}
}

// Edges of the Columns.
func (Columns) Edges() []ent.Edge {
	return nil
}
