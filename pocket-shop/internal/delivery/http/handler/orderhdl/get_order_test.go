package orderhdl_test

import (
	"errors"
	"net/http"
	"net/http/httptest"

	"pocket-shop/config"
	"pocket-shop/internal/core/order"
	"pocket-shop/internal/core/order/domain"
	"pocket-shop/internal/delivery/http/handler/orderhdl"

	mockorder "pocket-shop/mock/core/order"

	"github.com/gofiber/fiber/v3"
	"github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
)

const validOrderID = "550e8400-e29b-41d4-a716-446655440000"

var _ = ginkgo.Describe("GetOrder handler", func() {
	var (
		app     *fiber.App
		handler *orderhdl.Handler
		svc     *mockorder.MockOrderService
		cfg     *config.Config
	)

	ginkgo.BeforeEach(func() {
		cfg = &config.Config{RefSource: "test"}
		svc = mockorder.NewMockOrderService(ginkgo.GinkgoT())
		handler = orderhdl.NewHandler(svc, cfg)
		app = fiber.New()
		api := app.Group("/api/v1")
		handler.RegisterRoutes(api)
	})

	ginkgo.When("GetOrder succeeds", func() {
		ginkgo.It("returns 200 with id, status and redeemCode", func() {
			code := "CODE123"
			expected := &domain.Order{
				ID: validOrderID, RefID: "ref-1", RefSource: "test",
				Status: domain.StatusCompleted,
			}
			svc.EXPECT().GetOrder(mock.Anything, validOrderID).Return(expected, &code, nil)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/orders/"+validOrderID, nil)
			resp, err := app.Test(req, fiber.TestConfig{})
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(resp.Header.Get("Content-Type")).To(ContainSubstring("application/json"))
		})

		ginkgo.It("returns 200 with redeemCode null when service returns nil code", func() {
			expected := &domain.Order{
				ID: validOrderID, RefID: "ref-1", RefSource: "test",
				Status: domain.StatusProcessing,
			}
			svc.EXPECT().GetOrder(mock.Anything, validOrderID).Return(expected, nil, nil)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/orders/"+validOrderID, nil)
			resp, err := app.Test(req, fiber.TestConfig{})
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
		})
	})

	ginkgo.When("order not found", func() {
		ginkgo.It("returns 404 when service returns ErrOrderNotFound", func() {
			svc.EXPECT().GetOrder(mock.Anything, validOrderID).Return(nil, nil, order.ErrOrderNotFound)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/orders/"+validOrderID, nil)
			resp, err := app.Test(req, fiber.TestConfig{})
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
		})
	})

	ginkgo.When("GetOrder returns other error", func() {
		ginkgo.It("returns 500 with error response", func() {
			svc.EXPECT().GetOrder(mock.Anything, validOrderID).Return(nil, nil, errors.New("db error"))

			req := httptest.NewRequest(http.MethodGet, "/api/v1/orders/"+validOrderID, nil)
			resp, err := app.Test(req, fiber.TestConfig{})
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(http.StatusInternalServerError))
		})
	})

	ginkgo.When("id is invalid", func() {
		ginkgo.It("returns 400 when id is not a valid UUID", func() {
			req := httptest.NewRequest(http.MethodGet, "/api/v1/orders/not-a-uuid", nil)
			resp, err := app.Test(req, fiber.TestConfig{})
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
		})
	})
})
