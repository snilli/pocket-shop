package ordersvc

import (
	"context"
	"time"

	"pocket-shop/internal/core/order/domain"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func (s *Service) CreateOrder(ctx context.Context, refSource string) (*domain.Order, error) {
	txID, err := s.orderReservationRepo.Pull(ctx, s.cfg.RefSource)
	if err != nil {
		return nil, err
	}
	if txID != "" {
		o := domain.Create(txID, s.cfg.RefSource)
		o.Complete()
		created, err := s.repo.Create(ctx, *o)
		if err != nil {
			return nil, err
		}
		return created, nil
	}

	orderID := uuid.New().String()
	ezOrder, err := s.ez.CreateInstantOrder(ctx, orderID)
	if err != nil {
		return nil, err
	}
	o := domain.Create(ezOrder.TransactionId, s.cfg.RefSource)
	o.ID = orderID
	if ezOrder.GetStatus() == "COMPLETED" {
		code, codeErr := s.ez.GetFirstRedeemCode(ctx, ezOrder.TransactionId)
		if codeErr == nil && code != "" {
			o.Complete()
		}
	}
	created, err := s.repo.Create(ctx, *o)
	if err != nil {
		return nil, err
	}

	if created.Status == domain.StatusCompleted {
		return created, nil
	}

	timeout := time.Duration(s.cfg.OrderFulfillmentTimeoutSec) * time.Second
	pollInterval := time.Duration(s.cfg.PollIntervalSec) * time.Second
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		ezOrder, err := s.ez.GetOrder(ctx, o.RefID)
		if err != nil {
			log.Ctx(ctx).Warn().Err(err).Str("ref id", o.RefID).Msg("poll EZ status")
			time.Sleep(pollInterval)
			continue
		}
		status := ezOrder.GetStatus()
		switch status {
		case "COMPLETED":
			code, codeErr := s.ez.GetFirstRedeemCode(ctx, o.RefID)
			if codeErr == nil && code != "" {
				created.Complete()
				_ = s.repo.Save(ctx, created)
				return created, nil
			}
		case "CANCELLED":
			created.Cancel()
			_ = s.repo.Save(ctx, created)
			return created, nil
		}
		time.Sleep(pollInterval)
	}

	created.Cancel()
	_ = s.repo.Save(ctx, created)
	return created, nil
}
