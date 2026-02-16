package orderrepo

import (
	"context"
	"time"

	"pocket-shop/internal/core/order/domain"

	"github.com/google/uuid"
)

func (r *Repository) Create(ctx context.Context, o domain.Order) (*domain.Order, error) {
	now := time.Now()
	m := toModel(o)
	create := r.client.Order.Create().
		SetRefID(m.RefID).
		SetRefSource(m.RefSource).
		SetStatus(m.Status).
		SetCreatedAt(now)
	if o.ID != "" {
		id, err := uuid.Parse(o.ID)
		if err != nil {
			return nil, err
		}
		create.SetID(id)
	}
	if m.UpdatedAt != nil {
		create.SetUpdatedAt(*m.UpdatedAt)
	}
	if o.Status == domain.StatusCompleted {
		create.SetUsedAt(now)
	}

	entOrder, err := create.Save(ctx)
	if err != nil {
		return nil, err
	}

	return toDomain(entOrder), nil
}
