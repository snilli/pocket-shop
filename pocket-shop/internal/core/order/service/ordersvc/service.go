package ordersvc

import (
	"pocket-shop/config"
	"pocket-shop/internal/core/order"
	"pocket-shop/internal/port"
)

type Service struct {
	repo                 order.OrderRepository
	orderReservationRepo order.OrderReservationRepository
	ez                   port.EZClient
	cfg                  *config.Config
}

func New(repo order.OrderRepository, orderReservationRepo order.OrderReservationRepository, ez port.EZClient, cfg *config.Config) order.OrderService {
	return &Service{repo: repo, orderReservationRepo: orderReservationRepo, ez: ez, cfg: cfg}
}
