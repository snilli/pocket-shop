package port

import "context"

type EZClient interface {
	CreateInstantOrder(ctx context.Context, clientOrderID string) (*Order, error)
	GetOrder(ctx context.Context, refID string) (*Order, error)
	GetFirstRedeemCode(ctx context.Context, refID string) (string, error)
}
