package orderrepo

import (
	"context"
	"orm/ent/order"
	"time"
)

func (r *Repository) MarkRefIDUsed(ctx context.Context, refSource string, refIDs []string) error {
	if len(refIDs) == 0 {
		return nil
	}
	now := time.Now()
	_, err := r.client.Order.Update().
		Where(
			order.RefSourceEQ(refSource),
			order.UsedAtIsNil(),
			order.RefIDIn(refIDs...),
		).
		SetUsedAt(now).
		Save(ctx)
	return err
}
