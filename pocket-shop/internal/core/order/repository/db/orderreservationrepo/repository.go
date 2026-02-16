package orderreservationrepo

import (
	"orm"

	"pocket-shop/internal/core/order"
)

type Repository struct {
	client *orm.Client
}

func New(client *orm.Client) order.OrderReservationRepository {
	return &Repository{client: client}
}
