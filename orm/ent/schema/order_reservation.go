package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// OrderReservation holds available EZ transaction IDs for reuse (one-to-one with completed orders).
type OrderReservation struct {
	ent.Schema
}

func (OrderReservation) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "order_reservations"},
	}
}

func (OrderReservation) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").
			Positive(),
		field.String("ref_id").
			NotEmpty(),
		field.String("ref_source").
			NotEmpty(),
		field.Time("created_at").
			Optional().
			Nillable(),
		field.Time("used_at").
			Optional().
			Nillable().
			Comment("When it was taken/used; NULL = available"),
	}
}

func (OrderReservation) Edges() []ent.Edge {
	return nil
}

func (OrderReservation) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("ref_source", "ref_id").Unique(),
		index.Fields("ref_source", "used_at"),
	}
}
