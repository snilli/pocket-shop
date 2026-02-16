package orderreservationrepo

import (
	"context"
	"orm/ent"
	"orm/ent/orderreservation"
	"time"

	"entgo.io/ent/dialect/sql"
)

func (r *Repository) Pull(ctx context.Context, refSource string) (string, error) {
	tx, err := r.client.Tx(ctx)
	if err != nil {
		return "", err
	}
	defer func() { _ = tx.Rollback() }()

	node, err := tx.OrderReservation.Query().
		Where(
			orderreservation.RefSourceEQ(refSource),
			orderreservation.UsedAtIsNil(),
		).
		Order(orderreservation.ByID()).
		ForUpdate(sql.WithLockAction(sql.SkipLocked)).
		First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			_ = tx.Commit()
			return "", nil
		}
		return "", err
	}
	if err := tx.OrderReservation.UpdateOneID(node.ID).SetUsedAt(time.Now()).Exec(ctx); err != nil {
		return "", err
	}
	if err := tx.Commit(); err != nil {
		return "", err
	}
	return node.RefID, nil
}
