package orderhdl

import (
	"errors"

	"pocket-shop/internal/core/order"
	"pocket-shop/internal/delivery/http/response"
	"pocket-shop/internal/delivery/http/validation"

	"github.com/gofiber/fiber/v3"
)

// GetOrderResponse is the response for GET /orders/:id (assignment: id, status, redeemCode string | null).
type GetOrderResponse struct {
	ID         string  `json:"id"`
	Status     string  `json:"status"`
	RedeemCode *string `json:"redeemCode"` // null when PROCESSING or CANCELLED; required when COMPLETED
}

// GetOrder godoc
// @Summary Get order status and redeem code
// @Description Get order by local order ID. Returns status and redeemCode (null if PROCESSING or CANCELLED).
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Local Order ID"
// @Success 200 {object} GetOrderResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /orders/{id} [get]
func (h *Handler) GetOrder(c fiber.Ctx) error {
	params := GetOrderParams{ID: c.Params("id")}
	if err := validate.Struct(params); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(validation.ErrorResponse(err))
	}

	o, redeemCode, err := h.orderService.GetOrder(c.RequestCtx(), params.ID)
	if err != nil {
		if errors.Is(err, order.ErrOrderNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(response.ErrorResponse{
				Error:   "not_found",
				Message: "order not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
			Error:   "internal_error",
			Message: err.Error(),
		})
	}
	resp := GetOrderResponse{
		ID:     o.ID,
		Status: string(o.Status),
	}
	if redeemCode != nil {
		resp.RedeemCode = redeemCode
	}
	return c.JSON(resp)
}
