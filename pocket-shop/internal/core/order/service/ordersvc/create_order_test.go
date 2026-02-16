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

var _ = ginkgo.Describe("CreateOrder", func() {
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
			RefSource:                  "test",
			OrderFulfillmentTimeoutSec: 60,
			PollIntervalSec:            1,
		}
		repo = mockorder.NewMockOrderRepository(ginkgo.GinkgoT())
		reserv = mockorder.NewMockOrderReservationRepository(ginkgo.GinkgoT())
		ez = mockport.NewMockEZClient(ginkgo.GinkgoT())
		svc = ordersvc.New(repo, reserv, ez, cfg)
	})

	ginkgo.When("pool returns a txID", func() {
		ginkgo.It("returns completed order from pool without calling EZ", func() {
			txID := "pool-tx-1"
			expectedCreated := &domain.Order{
				ID:        "db-id-1",
				RefID:     txID,
				RefSource: cfg.RefSource,
				Status:    domain.StatusCompleted,
				CreatedAt: time.Now(),
			}
			reserv.EXPECT().Pull(ctx, cfg.RefSource).Return(txID, nil)
			repo.EXPECT().Create(ctx, mock.Anything).Return(expectedCreated, nil)

			created, err := svc.CreateOrder(ctx, cfg.RefSource)
			Expect(err).To(BeNil())
			Expect(created).NotTo(BeNil())
			Expect(created.RefID).To(Equal(txID))
			Expect(created.Status).To(Equal(domain.StatusCompleted))
		})

		ginkgo.It("returns error when repo.Create fails after pool", func() {
			reserv.EXPECT().Pull(ctx, cfg.RefSource).Return("pool-tx-1", nil)
			repo.EXPECT().Create(ctx, mock.Anything).Return(nil, errors.New("db error"))

			created, err := svc.CreateOrder(ctx, cfg.RefSource)
			Expect(err).NotTo(BeNil())
			Expect(created).To(BeNil())
		})
	})

	ginkgo.When("pool returns empty and EZ CreateInstantOrder succeeds", func() {
		ginkgo.It("creates order and returns when EZ returns COMPLETED with code", func() {
			reserv.EXPECT().Pull(ctx, cfg.RefSource).Return("", nil)

			orderID := "our-order-uuid"
			ezOrder := &port.Order{
				TransactionId:     "ez-tx-1",
				ClientOrderNumber: orderID,
				Status:            port.COMPLETED_OrderStatus,
			}
			ez.EXPECT().CreateInstantOrder(mock.Anything, mock.Anything).Return(ezOrder, nil)
			ez.EXPECT().GetFirstRedeemCode(mock.Anything, mock.Anything).Return("CODE123", nil)

			createdFromRepo := &domain.Order{
				ID:        orderID,
				RefID:     orderID,
				RefSource: cfg.RefSource,
				Status:    domain.StatusCompleted,
				CreatedAt: time.Now(),
			}

			repo.EXPECT().Create(ctx, mock.Anything).Return(createdFromRepo, nil)

			created, err := svc.CreateOrder(ctx, cfg.RefSource)
			Expect(err).To(BeNil())
			Expect(created).NotTo(BeNil())
			Expect(created.Status).To(Equal(domain.StatusCompleted))
			Expect(created.RefID).To(Equal(orderID))
		})

		ginkgo.It("returns error when CreateInstantOrder fails", func() {
			reserv.EXPECT().Pull(ctx, cfg.RefSource).Return("", nil)
			ez.EXPECT().CreateInstantOrder(mock.Anything, mock.Anything).Return(nil, errors.New("EZ API error"))

			created, err := svc.CreateOrder(ctx, cfg.RefSource)
			Expect(err).NotTo(BeNil())
			Expect(created).To(BeNil())
		})

		ginkgo.It("returns error when repo.Create fails after EZ success", func() {
			reserv.EXPECT().Pull(ctx, cfg.RefSource).Return("", nil)
			ezOrder := &port.Order{ClientOrderNumber: "oid", TransactionId: "tx1", Status: port.PROCESSING_OrderStatus}
			ez.EXPECT().CreateInstantOrder(mock.Anything, mock.Anything).Return(ezOrder, nil)
			repo.EXPECT().Create(ctx, mock.Anything).Return(nil, errors.New("db error"))

			created, err := svc.CreateOrder(ctx, cfg.RefSource)
			Expect(err).NotTo(BeNil())
			Expect(created).To(BeNil())
		})

		ginkgo.It("enters polling then cancels when EZ returns PROCESSING and timeout is 0", func() {
			reserv.EXPECT().Pull(ctx, cfg.RefSource).Return("", nil)
			ezOrder := &port.Order{ClientOrderNumber: "oid", TransactionId: "tx1", Status: port.PROCESSING_OrderStatus}
			ez.EXPECT().CreateInstantOrder(mock.Anything, mock.Anything).Return(ezOrder, nil)
			createdFromRepo := &domain.Order{
				ID: "oid", RefID: "oid", RefSource: cfg.RefSource,
				Status: domain.StatusProcessing, CreatedAt: time.Now(),
			}
			repo.EXPECT().Create(ctx, mock.Anything).Return(createdFromRepo, nil)
			repo.EXPECT().Save(ctx, mock.Anything).Return(nil)
			// timeout 0: loop never runs, we cancel and return
			cfg.OrderFulfillmentTimeoutSec = 0
			cfg.PollIntervalSec = 1

			created, err := svc.CreateOrder(ctx, cfg.RefSource)
			Expect(err).To(BeNil())
			Expect(created).NotTo(BeNil())
			Expect(created.Status).To(Equal(domain.StatusCancelled))
		})

		ginkgo.It("polls and returns when EZ returns CANCELLED", func() {
			reserv.EXPECT().Pull(ctx, cfg.RefSource).Return("", nil)
			ezOrder := &port.Order{ClientOrderNumber: "oid", TransactionId: "tx1", Status: port.PROCESSING_OrderStatus}
			ez.EXPECT().CreateInstantOrder(mock.Anything, mock.Anything).Return(ezOrder, nil)
			createdFromRepo := &domain.Order{
				ID: "oid", RefID: "oid", RefSource: cfg.RefSource,
				Status: domain.StatusProcessing, CreatedAt: time.Now(),
			}
			repo.EXPECT().Create(ctx, mock.Anything).Return(createdFromRepo, nil)
			cfg.OrderFulfillmentTimeoutSec = 60
			cfg.PollIntervalSec = 1
			ez.EXPECT().GetOrder(ctx, "oid").Return(&port.Order{TransactionId: "oid", Status: port.CANCELLED_OrderStatus}, nil)
			repo.EXPECT().Save(ctx, mock.Anything).Return(nil)

			created, err := svc.CreateOrder(ctx, cfg.RefSource)
			Expect(err).To(BeNil())
			Expect(created).NotTo(BeNil())
			Expect(created.Status).To(Equal(domain.StatusCancelled))
		})

		ginkgo.It("polls and returns completed when EZ returns COMPLETED with code", func() {
			reserv.EXPECT().Pull(ctx, cfg.RefSource).Return("", nil)
			ezOrder := &port.Order{ClientOrderNumber: "oid", TransactionId: "tx1", Status: port.PROCESSING_OrderStatus}
			ez.EXPECT().CreateInstantOrder(mock.Anything, mock.Anything).Return(ezOrder, nil)
			createdFromRepo := &domain.Order{
				ID: "oid", RefID: "oid", RefSource: cfg.RefSource,
				Status: domain.StatusProcessing, CreatedAt: time.Now(),
			}
			repo.EXPECT().Create(ctx, mock.Anything).Return(createdFromRepo, nil)
			cfg.OrderFulfillmentTimeoutSec = 60
			cfg.PollIntervalSec = 1
			ez.EXPECT().GetOrder(ctx, "oid").Return(&port.Order{TransactionId: "oid", Status: port.COMPLETED_OrderStatus}, nil)
			ez.EXPECT().GetFirstRedeemCode(ctx, "oid").Return("CODE99", nil)
			repo.EXPECT().Save(ctx, mock.Anything).Return(nil)

			created, err := svc.CreateOrder(ctx, cfg.RefSource)
			Expect(err).To(BeNil())
			Expect(created).NotTo(BeNil())
			Expect(created.Status).To(Equal(domain.StatusCompleted))
		})
	})

	ginkgo.When("Pull returns error", func() {
		ginkgo.It("returns error and does not call EZ or repo", func() {
			reserv.EXPECT().Pull(ctx, cfg.RefSource).Return("", errors.New("pool error"))

			created, err := svc.CreateOrder(ctx, cfg.RefSource)
			Expect(err).NotTo(BeNil())
			Expect(created).To(BeNil())
		})
	})
})
