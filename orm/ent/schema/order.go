package schema

import (
	"github.com/google/uuid"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type Order struct {
	ent.Schema
}

func (Order) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.String("ref_id").
			NotEmpty(),
		field.String("ref_source").
			NotEmpty(),
		field.Enum("status").
			Values("PROCESSING", "COMPLETED", "CANCELLED"),
		field.Time("created_at").
			Optional(),
		field.Time("updated_at").
			Optional().
			Nillable(),
		field.Time("used_at").
			Optional().
			Nillable(),
	}
}

func (Order) Edges() []ent.Edge {
	return nil
}

func (Order) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("ref_source", "ref_id", "status"),
		index.Fields("status"),
	}
}
