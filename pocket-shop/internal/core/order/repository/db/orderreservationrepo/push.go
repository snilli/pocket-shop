package orderreservationrepo

import (
	"context"
	"orm/ent"
	"time"
)

func (r *Repository) Push(ctx context.Context, refSource, refID string) error {
	err := r.client.OrderReservation.Create().
		SetRefID(refID).
		SetRefSource(refSource).
		SetCreatedAt(time.Now()).
		Exec(ctx)
	if err != nil && ent.IsConstraintError(err) {
		return nil
	}
	return err
}
