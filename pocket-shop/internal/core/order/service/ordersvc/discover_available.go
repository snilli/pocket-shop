package ordersvc

import (
	"context"
	"sync"

	"github.com/rs/zerolog/log"

	"pocket-shop/internal/core/order"
	"pocket-shop/internal/core/order/domain"

	"golang.org/x/sync/errgroup"
)

func (s *Service) DiscoverAvailable(ctx context.Context, refSource string) error {
	orders, err := s.repo.GetRecovery(ctx, refSource)
	if err != nil {
		return err
	}
	if len(orders) == 0 {
		return nil
	}

	var (
		mu       sync.Mutex
		pushedID []string
	)
	g, gCtx := errgroup.WithContext(ctx)
	for _, o := range orders {
		o := o
		g.Go(func() error {
			count, err := s.repo.Count(gCtx, order.OrderGetInput{
				RefSource: o.RefSource,
				RefID:     o.RefID,
				Status:    domain.StatusCompleted,
			})
			if err != nil || count > 0 {
				return nil
			}

			ezOrder, err := s.ez.GetOrder(gCtx, o.RefID)
			if err != nil || ezOrder == nil {
				return nil
			}

			code, err := s.ez.GetFirstRedeemCode(gCtx, o.RefID)
			if err != nil || code == "" {
				return nil
			}

			if err := s.orderReservationRepo.Push(gCtx, o.RefSource, o.RefID); err != nil {
				return nil
			}
			mu.Lock()
			pushedID = append(pushedID, o.RefID)
			mu.Unlock()
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

	allRefIDs := make([]string, 0, len(orders))
	for _, o := range orders {
		allRefIDs = append(allRefIDs, o.RefID)
	}
	if err := s.repo.MarkRefIDUsed(ctx, refSource, allRefIDs); err != nil {
		log.Warn().Err(err).Strs("ref_ids", allRefIDs).Msg("MarkRefIDUsed failed")
	}

	if len(pushedID) > 0 {
		log.Info().
			Str("ref_source", refSource).
			Int("recovered", len(pushedID)).
			Strs("ref_ids", pushedID).
			Msg("discover recovery")
	}

	return nil
}
