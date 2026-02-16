package orderhdl_test

import (
	"errors"
	"net/http"
	"net/http/httptest"

	"pocket-shop/config"
	"pocket-shop/internal/core/order/domain"
	"pocket-shop/internal/delivery/http/handler/orderhdl"

	mockorder "pocket-shop/mock/core/order"

	"github.com/gofiber/fiber/v3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
)

var _ = Describe("CreateOrder handler", func() {
	var (
		app     *fiber.App
		handler *orderhdl.Handler
		svc     *mockorder.MockOrderService
		cfg     *config.Config
	)

	BeforeEach(func() {
		cfg = &config.Config{RefSource: "test"}
		svc = mockorder.NewMockOrderService(GinkgoT())
		handler = orderhdl.NewHandler(svc, cfg)
		app = fiber.New()
		api := app.Group("/api/v1")
		handler.RegisterRoutes(api)
	})

	When("CreateOrder succeeds", func() {
		It("returns 200 with id and status", func() {
			expected := &domain.Order{
				ID: "order-123", RefID: "ref-1", RefSource: "test",
				Status: domain.StatusCompleted,
			}
			svc.EXPECT().CreateOrder(mock.Anything, "test").Return(expected, nil)

			req := httptest.NewRequest(http.MethodPost, "/api/v1/orders", nil)
			resp, err := app.Test(req, fiber.TestConfig{})
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(resp.Header.Get("Content-Type")).To(ContainSubstring("application/json"))
		})
	})

	When("CreateOrder returns error", func() {
		It("returns 500 with error response", func() {
			svc.EXPECT().CreateOrder(mock.Anything, "test").Return(nil, errors.New("service error"))

			req := httptest.NewRequest(http.MethodPost, "/api/v1/orders", nil)
			resp, err := app.Test(req, fiber.TestConfig{})
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(http.StatusInternalServerError))
		})
	})
})
