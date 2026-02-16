package order

import (
	"context"

	"pocket-shop/internal/core/order/domain"
)

type OrderGetInput struct {
	RefSource string
	RefID     string
	Status    domain.OrderStatus
}

type OrderRepository interface {
	Create(ctx context.Context, o domain.Order) (*domain.Order, error)
	GetByID(ctx context.Context, id string) (*domain.Order, error)
	Save(ctx context.Context, o *domain.Order) error
	GetRecovery(ctx context.Context, refSource string) ([]domain.Order, error)
	Count(ctx context.Context, input OrderGetInput) (int, error)
	MarkRefIDUsed(ctx context.Context, refSource string, refIDs []string) error
}

type OrderReservationRepository interface {
	Pull(ctx context.Context, refSource string) (string, error)
	Push(ctx context.Context, refSource, refID string) error
}
