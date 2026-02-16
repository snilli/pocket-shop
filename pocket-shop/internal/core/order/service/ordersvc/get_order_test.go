package ordersvc_test

import (
	"context"
	"errors"
	"time"

	"pocket-shop/config"
	"pocket-shop/internal/core/order"
	"pocket-shop/internal/core/order/domain"
	"pocket-shop/internal/core/order/service/ordersvc"
	"pocket-shop/internal/port"

	"github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	mockorder "pocket-shop/mock/core/order"
	mockport "pocket-shop/mock/port"
)

var _ = ginkgo.Describe("GetOrder", func() {
	var (
		ctx    context.Context
		cfg    *config.Config
		repo   *mockorder.MockOrderRepository
		reserv *mockorder.MockOrderReservationRepository
		ez     *mockport.MockEZClient
		svc    order.OrderService
	)

	ginkgo.BeforeEach(func() {
		ctx = context.Background()
		cfg = &config.Config{
			RefSource:                    "test",
			OrderFulfillmentTimeoutSec:   60,
		}
		repo = mockorder.NewMockOrderRepository(ginkgo.GinkgoT())
		reserv = mockorder.NewMockOrderReservationRepository(ginkgo.GinkgoT())
		ez = mockport.NewMockEZClient(ginkgo.GinkgoT())
		svc = ordersvc.New(repo, reserv, ez, cfg)
	})

	ginkgo.When("order not found", func() {
		ginkgo.It("returns ErrOrderNotFound when repo.GetByID returns nil", func() {
			repo.EXPECT().GetByID(ctx, "id-1").Return(nil, nil)

			o, code, err := svc.GetOrder(ctx, "id-1")
			Expect(err).To(Equal(order.ErrOrderNotFound))
			Expect(o).To(BeNil())
			Expect(code).To(BeNil())
		})

		ginkgo.It("returns error when repo.GetByID fails", func() {
			repo.EXPECT().GetByID(ctx, "id-1").Return(nil, errors.New("db error"))

			o, code, err := svc.GetOrder(ctx, "id-1")
			Expect(err).NotTo(BeNil())
			Expect(o).To(BeNil())
			Expect(code).To(BeNil())
		})
	})

	ginkgo.When("order is cancelled", func() {
		ginkgo.It("returns order with nil code without calling EZ", func() {
			expected := &domain.Order{
				ID: "id-1", RefID: "ref-1", RefSource: "test",
				Status: domain.StatusCancelled, CreatedAt: time.Now(),
			}
			repo.EXPECT().GetByID(ctx, "id-1").Return(expected, nil)

			o, code, err := svc.GetOrder(ctx, "id-1")
			Expect(err).To(BeNil())
			Expect(o).To(Equal(expected))
			Expect(code).To(BeNil())
		})
	})

	ginkgo.When("order is completed", func() {
		ginkgo.It("returns order and redeem code when GetFirstRedeemCode returns code", func() {
			expected := &domain.Order{
				ID: "id-1", RefID: "ref-1", RefSource: "test",
				Status: domain.StatusCompleted, CreatedAt: time.Now(),
			}
			repo.EXPECT().GetByID(ctx, "id-1").Return(expected, nil)
			ez.EXPECT().GetFirstRedeemCode(ctx, "ref-1").Return("CODE123", nil)

			o, code, err := svc.GetOrder(ctx, "id-1")
			Expect(err).To(BeNil())
			Expect(o).To(Equal(expected))
			Expect(code).NotTo(BeNil())
			Expect(*code).To(Equal("CODE123"))
		})

		ginkgo.It("returns order and nil code when GetFirstRedeemCode returns empty", func() {
			expected := &domain.Order{
				ID: "id-1", RefID: "ref-1", RefSource: "test",
				Status: domain.StatusCompleted, CreatedAt: time.Now(),
			}
			repo.EXPECT().GetByID(ctx, "id-1").Return(expected, nil)
			ez.EXPECT().GetFirstRedeemCode(ctx, "ref-1").Return("", nil)

			o, code, err := svc.GetOrder(ctx, "id-1")
			Expect(err).To(BeNil())
			Expect(o).To(Equal(expected))
			Expect(code).To(BeNil())
		})

		ginkgo.It("returns error when GetFirstRedeemCode fails for completed order", func() {
			expected := &domain.Order{
				ID: "id-1", RefID: "ref-1", RefSource: "test",
				Status: domain.StatusCompleted, CreatedAt: time.Now(),
			}
			repo.EXPECT().GetByID(ctx, "id-1").Return(expected, nil)
			ez.EXPECT().GetFirstRedeemCode(ctx, "ref-1").Return("", errors.New("ez error"))

			o, code, err := svc.GetOrder(ctx, "id-1")
			Expect(err).NotTo(BeNil())
			Expect(o).To(Equal(expected))
			Expect(code).To(BeNil())
		})
	})

	ginkgo.When("order is processing", func() {
		ginkgo.It("returns order with completed status and code when EZ returns COMPLETED with code", func() {
			expected := &domain.Order{
				ID: "id-1", RefID: "ref-1", RefSource: "test",
				Status: domain.StatusProcessing, CreatedAt: time.Now(),
			}
			repo.EXPECT().GetByID(ctx, "id-1").Return(expected, nil)
			ez.EXPECT().GetOrder(ctx, "ref-1").Return(&port.Order{
				TransactionId: "ref-1", ClientOrderNumber: "c1", Status: port.COMPLETED_OrderStatus,
			}, nil)
			ez.EXPECT().GetFirstRedeemCode(ctx, "ref-1").Return("CODE456", nil)
			repo.EXPECT().Save(ctx, mock.Anything).Return(nil)

			o, code, err := svc.GetOrder(ctx, "id-1")
			Expect(err).To(BeNil())
			Expect(o).NotTo(BeNil())
			Expect(o.Status).To(Equal(domain.StatusCompleted))
			Expect(code).NotTo(BeNil())
			Expect(*code).To(Equal("CODE456"))
		})

		ginkgo.It("returns order cancelled when past fulfillment timeout", func() {
			createdAt := time.Now().Add(-2 * time.Hour)
			expected := &domain.Order{
				ID: "id-1", RefID: "ref-1", RefSource: "test",
				Status: domain.StatusProcessing, CreatedAt: createdAt,
			}
			repo.EXPECT().GetByID(ctx, "id-1").Return(expected, nil)
			repo.EXPECT().Save(ctx, mock.Anything).Return(nil)

			o, code, err := svc.GetOrder(ctx, "id-1")
			Expect(err).To(BeNil())
			Expect(o).NotTo(BeNil())
			Expect(o.Status).To(Equal(domain.StatusCancelled))
			Expect(code).To(BeNil())
		})

		ginkgo.It("returns order with nil code when GetOrder returns error", func() {
			expected := &domain.Order{
				ID: "id-1", RefID: "ref-1", RefSource: "test",
				Status: domain.StatusProcessing, CreatedAt: time.Now(),
			}
			repo.EXPECT().GetByID(ctx, "id-1").Return(expected, nil)
			ez.EXPECT().GetOrder(ctx, "ref-1").Return(nil, errors.New("ez error"))

			o, code, err := svc.GetOrder(ctx, "id-1")
			Expect(err).To(BeNil())
			Expect(o).To(Equal(expected))
			Expect(code).To(BeNil())
		})

		ginkgo.It("returns order with nil code when GetOrder returns nil", func() {
			expected := &domain.Order{
				ID: "id-1", RefID: "ref-1", RefSource: "test",
				Status: domain.StatusProcessing, CreatedAt: time.Now(),
			}
			repo.EXPECT().GetByID(ctx, "id-1").Return(expected, nil)
			ez.EXPECT().GetOrder(ctx, "ref-1").Return(nil, nil)

			o, code, err := svc.GetOrder(ctx, "id-1")
			Expect(err).To(BeNil())
			Expect(o).To(Equal(expected))
			Expect(code).To(BeNil())
		})

		ginkgo.It("returns order with nil code when EZ returns non-COMPLETED", func() {
			expected := &domain.Order{
				ID: "id-1", RefID: "ref-1", RefSource: "test",
				Status: domain.StatusProcessing, CreatedAt: time.Now(),
			}
			repo.EXPECT().GetByID(ctx, "id-1").Return(expected, nil)
			ez.EXPECT().GetOrder(ctx, "ref-1").Return(&port.Order{
				TransactionId: "ref-1", Status: port.PROCESSING_OrderStatus,
			}, nil)

			o, code, err := svc.GetOrder(ctx, "id-1")
			Expect(err).To(BeNil())
			Expect(o).To(Equal(expected))
			Expect(code).To(BeNil())
		})

		ginkgo.It("returns order with nil code when EZ COMPLETED but GetFirstRedeemCode returns empty", func() {
			expected := &domain.Order{
				ID: "id-1", RefID: "ref-1", RefSource: "test",
				Status: domain.StatusProcessing, CreatedAt: time.Now(),
			}
			repo.EXPECT().GetByID(ctx, "id-1").Return(expected, nil)
			ez.EXPECT().GetOrder(ctx, "ref-1").Return(&port.Order{
				TransactionId: "ref-1", Status: port.COMPLETED_OrderStatus,
			}, nil)
			ez.EXPECT().GetFirstRedeemCode(ctx, "ref-1").Return("", nil)

			o, code, err := svc.GetOrder(ctx, "id-1")
			Expect(err).To(BeNil())
			Expect(o).NotTo(BeNil())
			Expect(code).To(BeNil())
		})
	})
})
