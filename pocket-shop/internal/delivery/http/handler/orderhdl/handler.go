package orderhdl

import (
	"pocket-shop/config"
	"pocket-shop/internal/core/order"

	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	orderService order.OrderService
	cfg          *config.Config
}

func NewHandler(orderService order.OrderService, cfg *config.Config) *Handler {
	return &Handler{orderService: orderService, cfg: cfg}
}

func (h *Handler) RegisterRoutes(router fiber.Router) {
	orders := router.Group("/orders")
	orders.Post("/", h.CreateOrder)
	orders.Get("/:id", h.GetOrder)
}
