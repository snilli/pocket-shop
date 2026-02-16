package orderhdl

import (
	"pocket-shop/internal/delivery/http/response"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog/log"
)

// CreateOrderResponse is the response for POST /orders (assignment: id, status).
type CreateOrderResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

// CreateOrder godoc
// @Summary Create order
// @Description Create a new instant gift card order. Reuses an available EZ order from pool if any; otherwise creates a new EZ order and polls until completed or timeout.
// @Tags orders
// @Accept json
// @Produce json
// @Success 200 {object} CreateOrderResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /orders [post]
func (h *Handler) CreateOrder(c fiber.Ctx) error {
	o, err := h.orderService.CreateOrder(c.RequestCtx(), h.cfg.RefSource)
	if err != nil {
		log.Error().Err(err).Msg("CreateOrder failed")
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
			Error:   "internal_error",
			Message: err.Error(),
		})
	}
	return c.JSON(CreateOrderResponse{
		ID:     o.ID,
		Status: string(o.Status),
	})
}
