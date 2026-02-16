package orderrepo

import (
	"context"
	"orm/ent"

	"pocket-shop/internal/core/order/domain"

	"github.com/google/uuid"
)

func (r *Repository) GetByID(ctx context.Context, ID string) (*domain.Order, error) {
	entOrder, err := r.client.Order.Get(ctx, uuid.MustParse(ID))
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return toDomain(entOrder), nil
}
