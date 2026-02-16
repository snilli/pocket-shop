package orderrepo

import (
	"orm"

	"pocket-shop/internal/core/order"
)

type Repository struct {
	client *orm.Client
}

func New(client *orm.Client) order.OrderRepository {
	return &Repository{client: client}
}
