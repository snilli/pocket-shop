package orderrepo

import (
	"context"
	"orm/ent/order"

	"pocket-shop/internal/core/order/domain"
)

func (r *Repository) GetRecovery(ctx context.Context, refSource string) ([]domain.Order, error) {
	candidateNodes, err := r.client.Order.Query().
		Where(
			order.RefSourceEQ(refSource),
			order.UsedAtIsNil(),
			order.Or(
				order.StatusEQ(order.StatusPROCESSING),
				order.StatusEQ(order.StatusCANCELLED),
			),
		).
		All(ctx)
	if err != nil {
		return nil, err
	}

	if len(candidateNodes) == 0 {
		return nil, nil
	}

	candidates := make([]string, 0, len(candidateNodes))
	orders := make([]domain.Order, 0, len(candidateNodes))
	for _, o := range candidateNodes {
		order := toDomain(o)
		orders = append(orders, *order)
		candidates = append(candidates, o.RefID)
	}

	usedNodes, err := r.client.Order.Query().
		Where(
			order.RefSourceEQ(refSource),
			order.StatusEQ(order.StatusCOMPLETED),
			order.RefIDIn(candidates...),
		).
		All(ctx)
	if err != nil {
		return nil, err
	}
	usedSet := make(map[string]bool)
	for _, n := range usedNodes {
		usedSet[n.RefID] = true
	}
	var out []domain.Order
	seen := make(map[string]bool)
	for _, o := range orders {
		if !usedSet[o.RefID] && !seen[o.RefID] {
			seen[o.RefID] = true
			out = append(out, o)
		}
	}
	return out, nil
}
