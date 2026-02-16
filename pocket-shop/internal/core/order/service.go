package order

import (
	"context"

	"pocket-shop/internal/core/order/domain"
)

type OrderService interface {
	CreateOrder(ctx context.Context, refSource string) (*domain.Order, error)
	GetOrder(ctx context.Context, id string) (*domain.Order, *string, error)
	DiscoverAvailable(ctx context.Context, refSource string) error
}
