package orderrepo

import (
	"context"
	entorder "orm/ent/order"

	"pocket-shop/internal/core/order"
)

func (r *Repository) Count(ctx context.Context, input order.OrderGetInput) (int, error) {
	n, err := r.client.Order.Query().
		Where(
			entorder.RefSourceEQ(input.RefSource),
			entorder.RefIDEQ(input.RefID),
			entorder.StatusEQ(toEntStatus(input.Status)),
		).
		Count(ctx)
	if err != nil {
		return 0, err
	}
	return n, nil
}
