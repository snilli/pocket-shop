package orderrepo

import (
	"context"
	"time"

	"github.com/google/uuid"

	"pocket-shop/internal/core/order/domain"
)

func (r *Repository) Save(ctx context.Context, o *domain.Order) error {
	if o == nil {
		return nil
	}
	id, err := uuid.Parse(o.ID)
	if err != nil {
		return err
	}
	upd := r.client.Order.UpdateOneID(id).
		SetRefID(o.RefID).
		SetRefSource(o.RefSource).
		SetStatus(toEntStatus(o.Status)).
		SetCreatedAt(o.CreatedAt)
	if o.UpdatedAt != nil {
		upd = upd.SetUpdatedAt(*o.UpdatedAt)
	}

	if o.Status == domain.StatusCompleted {
		upd = upd.SetUsedAt(time.Now())
	}

	_, err = upd.Save(ctx)
	return err
}
