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

	mockorder "pocket-shop/mock/core/order"
	mockport "pocket-shop/mock/port"

	"github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
)

var _ = ginkgo.Describe("DiscoverAvailable", func() {
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
		cfg = &config.Config{RefSource: "test"}
		repo = mockorder.NewMockOrderRepository(ginkgo.GinkgoT())
		reserv = mockorder.NewMockOrderReservationRepository(ginkgo.GinkgoT())
		ez = mockport.NewMockEZClient(ginkgo.GinkgoT())
		svc = ordersvc.New(repo, reserv, ez, cfg)
	})

	ginkgo.When("GetRecovery fails", func() {
		ginkgo.It("returns error", func() {
			repo.EXPECT().GetRecovery(ctx, "test").Return(nil, errors.New("db error"))

			err := svc.DiscoverAvailable(ctx, "test")
			Expect(err).NotTo(BeNil())
		})
	})

	ginkgo.When("no recovery orders", func() {
		ginkgo.It("returns nil without calling EZ or Push", func() {
			repo.EXPECT().GetRecovery(ctx, "test").Return([]domain.Order{}, nil)

			err := svc.DiscoverAvailable(ctx, "test")
			Expect(err).To(BeNil())
		})
	})

	ginkgo.When("recovery orders exist and one is completed with code", func() {
		ginkgo.It("pushes refID to pool and calls MarkRefIDUsed", func() {
			orders := []domain.Order{
				{ID: "o1", RefID: "ref-1", RefSource: "test", Status: domain.StatusProcessing, CreatedAt: time.Now()},
			}
			repo.EXPECT().GetRecovery(ctx, "test").Return(orders, nil)
			repo.EXPECT().Count(mock.Anything, mock.Anything).Return(0, nil)
			ez.EXPECT().GetOrder(mock.Anything, "ref-1").Return(&port.Order{
				TransactionId: "ref-1", Status: port.COMPLETED_OrderStatus,
			}, nil)
			ez.EXPECT().GetFirstRedeemCode(mock.Anything, "ref-1").Return("CODE1", nil)
			reserv.EXPECT().Push(mock.Anything, "test", "ref-1").Return(nil)
			repo.EXPECT().MarkRefIDUsed(ctx, "test", []string{"ref-1"}).Return(nil)

			err := svc.DiscoverAvailable(ctx, "test")
			Expect(err).To(BeNil())
		})

		ginkgo.It("continues when MarkRefIDUsed fails after push", func() {
			orders := []domain.Order{
				{ID: "o1", RefID: "ref-1", RefSource: "test", Status: domain.StatusProcessing, CreatedAt: time.Now()},
			}
			repo.EXPECT().GetRecovery(ctx, "test").Return(orders, nil)
			repo.EXPECT().Count(mock.Anything, mock.Anything).Return(0, nil)
			ez.EXPECT().GetOrder(mock.Anything, "ref-1").Return(&port.Order{
				TransactionId: "ref-1", Status: port.COMPLETED_OrderStatus,
			}, nil)
			ez.EXPECT().GetFirstRedeemCode(mock.Anything, "ref-1").Return("CODE1", nil)
			reserv.EXPECT().Push(mock.Anything, "test", "ref-1").Return(nil)
			repo.EXPECT().MarkRefIDUsed(ctx, "test", []string{"ref-1"}).Return(errors.New("mark failed"))

			err := svc.DiscoverAvailable(ctx, "test")
			Expect(err).To(BeNil())
		})
	})

	ginkgo.When("recovery orders exist but none get pushed", func() {
		ginkgo.It("returns nil when Count returns error so order is skipped", func() {
			orders := []domain.Order{
				{ID: "o1", RefID: "ref-1", RefSource: "test", Status: domain.StatusProcessing, CreatedAt: time.Now()},
			}
			repo.EXPECT().GetRecovery(ctx, "test").Return(orders, nil)
			repo.EXPECT().Count(mock.Anything, mock.Anything).Return(0, errors.New("count error"))

			err := svc.DiscoverAvailable(ctx, "test")
			Expect(err).To(BeNil())
		})

		ginkgo.It("returns nil when Count > 0 so order is skipped", func() {
			orders := []domain.Order{
				{ID: "o1", RefID: "ref-1", RefSource: "test", Status: domain.StatusProcessing, CreatedAt: time.Now()},
			}
			repo.EXPECT().GetRecovery(ctx, "test").Return(orders, nil)
			repo.EXPECT().Count(mock.Anything, mock.Anything).Return(1, nil)

			err := svc.DiscoverAvailable(ctx, "test")
			Expect(err).To(BeNil())
		})

		ginkgo.It("returns nil when GetOrder returns error", func() {
			orders := []domain.Order{
				{ID: "o1", RefID: "ref-1", RefSource: "test", Status: domain.StatusProcessing, CreatedAt: time.Now()},
			}
			repo.EXPECT().GetRecovery(ctx, "test").Return(orders, nil)
			repo.EXPECT().Count(mock.Anything, mock.Anything).Return(0, nil)
			ez.EXPECT().GetOrder(mock.Anything, "ref-1").Return(nil, errors.New("ez error"))

			err := svc.DiscoverAvailable(ctx, "test")
			Expect(err).To(BeNil())
		})

		ginkgo.It("returns nil when GetOrder returns nil order", func() {
			orders := []domain.Order{
				{ID: "o1", RefID: "ref-1", RefSource: "test", Status: domain.StatusProcessing, CreatedAt: time.Now()},
			}
			repo.EXPECT().GetRecovery(ctx, "test").Return(orders, nil)
			repo.EXPECT().Count(mock.Anything, mock.Anything).Return(0, nil)
			ez.EXPECT().GetOrder(mock.Anything, "ref-1").Return(nil, nil)

			err := svc.DiscoverAvailable(ctx, "test")
			Expect(err).To(BeNil())
		})

		ginkgo.It("returns nil when GetOrder returns non-COMPLETED", func() {
			orders := []domain.Order{
				{ID: "o1", RefID: "ref-1", RefSource: "test", Status: domain.StatusProcessing, CreatedAt: time.Now()},
			}
			repo.EXPECT().GetRecovery(ctx, "test").Return(orders, nil)
			repo.EXPECT().Count(mock.Anything, mock.Anything).Return(0, nil)
			ez.EXPECT().GetOrder(mock.Anything, "ref-1").Return(&port.Order{
				TransactionId: "ref-1", Status: port.PROCESSING_OrderStatus,
			}, nil)

			err := svc.DiscoverAvailable(ctx, "test")
			Expect(err).To(BeNil())
		})

		ginkgo.It("returns nil when GetFirstRedeemCode returns error", func() {
			orders := []domain.Order{
				{ID: "o1", RefID: "ref-1", RefSource: "test", Status: domain.StatusProcessing, CreatedAt: time.Now()},
			}
			repo.EXPECT().GetRecovery(ctx, "test").Return(orders, nil)
			repo.EXPECT().Count(mock.Anything, mock.Anything).Return(0, nil)
			ez.EXPECT().GetOrder(mock.Anything, "ref-1").Return(&port.Order{
				TransactionId: "ref-1", Status: port.COMPLETED_OrderStatus,
			}, nil)
			ez.EXPECT().GetFirstRedeemCode(mock.Anything, "ref-1").Return("", errors.New("code error"))

			err := svc.DiscoverAvailable(ctx, "test")
			Expect(err).To(BeNil())
		})

		ginkgo.It("returns nil when GetFirstRedeemCode returns empty", func() {
			orders := []domain.Order{
				{ID: "o1", RefID: "ref-1", RefSource: "test", Status: domain.StatusProcessing, CreatedAt: time.Now()},
			}
			repo.EXPECT().GetRecovery(ctx, "test").Return(orders, nil)
			repo.EXPECT().Count(mock.Anything, mock.Anything).Return(0, nil)
			ez.EXPECT().GetOrder(mock.Anything, "ref-1").Return(&port.Order{
				TransactionId: "ref-1", Status: port.COMPLETED_OrderStatus,
			}, nil)
			ez.EXPECT().GetFirstRedeemCode(mock.Anything, "ref-1").Return("", nil)

			err := svc.DiscoverAvailable(ctx, "test")
			Expect(err).To(BeNil())
		})

		ginkgo.It("returns nil when Push fails", func() {
			orders := []domain.Order{
				{ID: "o1", RefID: "ref-1", RefSource: "test", Status: domain.StatusProcessing, CreatedAt: time.Now()},
			}
			repo.EXPECT().GetRecovery(ctx, "test").Return(orders, nil)
			repo.EXPECT().Count(mock.Anything, mock.Anything).Return(0, nil)
			ez.EXPECT().GetOrder(mock.Anything, "ref-1").Return(&port.Order{
				TransactionId: "ref-1", Status: port.COMPLETED_OrderStatus,
			}, nil)
			ez.EXPECT().GetFirstRedeemCode(mock.Anything, "ref-1").Return("CODE1", nil)
			reserv.EXPECT().Push(mock.Anything, "test", "ref-1").Return(errors.New("push failed"))

			err := svc.DiscoverAvailable(ctx, "test")
			Expect(err).To(BeNil())
		})
	})
})
