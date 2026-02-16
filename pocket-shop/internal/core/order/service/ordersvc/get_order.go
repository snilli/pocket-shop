package ordersvc

import (
	"context"
	"time"

	"pocket-shop/internal/core/order"
	"pocket-shop/internal/core/order/domain"
)

func (s *Service) GetOrder(ctx context.Context, id string) (*domain.Order, *string, error) {
	o, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	if o == nil {
		return nil, nil, order.ErrOrderNotFound
	}

	switch o.Status {
	case domain.StatusCancelled:
		return o, nil, nil
	case domain.StatusCompleted:
		code, err := s.ez.GetFirstRedeemCode(ctx, o.RefID)
		if err != nil {
			return o, nil, err
		}
		if code != "" {
			return o, &code, nil
		}
		return o, nil, nil
	case domain.StatusProcessing:
		timeout := time.Duration(s.cfg.OrderFulfillmentTimeoutSec) * time.Second
		deadline := o.CreatedAt.Add(timeout)
		if time.Now().After(deadline) {
			o.Cancel()
			_ = s.repo.Save(ctx, o)
			return o, nil, nil
		}
		ezOrder, err := s.ez.GetOrder(ctx, o.RefID)
		if err != nil || ezOrder == nil {
			return o, nil, nil
		}
		if ezOrder.GetStatus() == "COMPLETED" {
			code, _ := s.ez.GetFirstRedeemCode(ctx, o.RefID)
			if code != "" {
				o.Complete()
				_ = s.repo.Save(ctx, o)
				return o, &code, nil
			}
		}
		return o, nil, nil
	}
	return o, nil, nil
}
