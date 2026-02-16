package orderrepo_test

import (
	"context"
	"time"

	"pocket-shop/internal/core/order"
	"pocket-shop/internal/core/order/domain"
	"pocket-shop/internal/core/order/repository/db/orderrepo"
	"pocket-shop/internal/testutil"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Order repository integration", Label("integration"), func() {
	var (
		ctx  context.Context
		pg   *testutil.PostgresContainer
		repo order.OrderRepository
	)

	BeforeEach(func() {
		ctx = context.Background()
		var err error
		pg, err = testutil.NewPostgresContainer(ctx)
		Expect(err).NotTo(HaveOccurred())
		repo = orderrepo.New(pg.Client)
	})

	AfterEach(func() {
		if pg != nil {
			Expect(pg.Close(ctx)).To(Succeed())
		}
	})

	DescribeTable("Create",
		func(setup func() domain.Order, expectID bool) {
			o := setup()
			created, err := repo.Create(ctx, o)
			Expect(err).NotTo(HaveOccurred())
			Expect(created).NotTo(BeNil())
			Expect(created.RefID).To(Equal(o.RefID))
			Expect(created.RefSource).To(Equal(o.RefSource))
			Expect(created.Status).To(Equal(o.Status))
			if expectID {
				Expect(created.ID).NotTo(BeEmpty())
			}
		},
		Entry("with explicit ID", func() domain.Order {
			id := uuid.New().String()
			o := domain.Create("ref-1", "source")
			o.ID = id
			return *o
		}, true),
		Entry("without ID (auto UUID)", func() domain.Order {
			o := domain.Create("ref-2", "source")
			return *o
		}, true),
		Entry("completed status", func() domain.Order {
			o := domain.Create("ref-3", "source")
			o.Complete()
			return *o
		}, true),
	)

	When("Create and GetByID", func() {
		It("saves and retrieves order by ID", func() {
			o := domain.Create("ref-get", "source")
			o.ID = uuid.New().String()
			created, err := repo.Create(ctx, *o)
			Expect(err).NotTo(HaveOccurred())
			Expect(created.ID).NotTo(BeEmpty())

			got, err := repo.GetByID(ctx, created.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(got).NotTo(BeNil())
			Expect(got.ID).To(Equal(created.ID))
			Expect(got.RefID).To(Equal("ref-get"))
			Expect(got.RefSource).To(Equal("source"))
		})

		It("returns nil nil when order not found", func() {
			got, err := repo.GetByID(ctx, uuid.New().String())
			Expect(err).NotTo(HaveOccurred())
			Expect(got).To(BeNil())
		})
	})

	When("Save", func() {
		It("updates order status", func() {
			o := domain.Create("ref-save", "source")
			o.ID = uuid.New().String()
			created, err := repo.Create(ctx, *o)
			Expect(err).NotTo(HaveOccurred())

			created.Complete()
			Expect(repo.Save(ctx, created)).To(Succeed())

			got, err := repo.GetByID(ctx, created.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(got.Status).To(Equal(domain.StatusCompleted))
		})
	})

	When("Count", func() {
		It("returns 0 when no matching orders", func() {
			n, err := repo.Count(ctx, order.OrderGetInput{
				RefSource: "src",
				RefID:     "r1",
				Status:    domain.StatusCompleted,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(n).To(Equal(0))
		})

		It("returns count of matching orders", func() {
			o := domain.Create("r1", "src")
			o.ID = uuid.New().String()
			o.Complete()
			_, err := repo.Create(ctx, *o)
			Expect(err).NotTo(HaveOccurred())
			o2 := domain.Create("r1", "src")
			o2.ID = uuid.New().String()
			o2.Complete()
			_, err = repo.Create(ctx, *o2)
			Expect(err).NotTo(HaveOccurred())

			n, err := repo.Count(ctx, order.OrderGetInput{
				RefSource: "src",
				RefID:     "r1",
				Status:    domain.StatusCompleted,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(n).To(Equal(2))
		})
	})

	When("GetRecovery", func() {
		It("returns empty when no recovery candidates", func() {
			list, err := repo.GetRecovery(ctx, "src")
			Expect(err).NotTo(HaveOccurred())
			Expect(list).To(BeEmpty())
		})

		It("returns processing orders not yet used (no completed with same ref_id)", func() {
			o := domain.Create("ref-recovery", "src")
			o.ID = uuid.New().String()
			o.Status = domain.StatusProcessing
			o.CreatedAt = time.Now()
			_, err := repo.Create(ctx, *o)
			Expect(err).NotTo(HaveOccurred())

			list, err := repo.GetRecovery(ctx, "src")
			Expect(err).NotTo(HaveOccurred())
			Expect(list).To(HaveLen(1))
			Expect(list[0].RefID).To(Equal("ref-recovery"))
		})

		It("excludes ref_ids that have a completed order", func() {
			// Processing with ref-recovery
			o1 := domain.Create("ref-recovery", "src")
			o1.ID = uuid.New().String()
			o1.Status = domain.StatusProcessing
			_, err := repo.Create(ctx, *o1)
			Expect(err).NotTo(HaveOccurred())
			// Completed with same ref_recovery -> used, so recovery should exclude it
			o2 := domain.Create("ref-recovery", "src")
			o2.ID = uuid.New().String()
			o2.Complete()
			_, err = repo.Create(ctx, *o2)
			Expect(err).NotTo(HaveOccurred())

			list, err := repo.GetRecovery(ctx, "src")
			Expect(err).NotTo(HaveOccurred())
			Expect(list).To(BeEmpty())
		})
	})

	When("MarkRefIDUsed", func() {
		It("sets used_at for matching orders", func() {
			o := domain.Create("ref-mark", "src")
			o.ID = uuid.New().String()
			o.Status = domain.StatusProcessing
			_, err := repo.Create(ctx, *o)
			Expect(err).NotTo(HaveOccurred())

			Expect(repo.MarkRefIDUsed(ctx, "src", []string{"ref-mark"})).To(Succeed())

			// GetRecovery excludes used ref_ids (completed with that ref_id or used_at set)
			// MarkRefIDUsed sets used_at, so same ref_id with used_at set - GetRecovery
			// filters by UsedAtIsNil(), so this order will no longer appear in recovery
			list, err := repo.GetRecovery(ctx, "src")
			Expect(err).NotTo(HaveOccurred())
			Expect(list).To(BeEmpty())
		})

		It("is no-op when refIDs empty", func() {
			Expect(repo.MarkRefIDUsed(ctx, "src", nil)).To(Succeed())
			Expect(repo.MarkRefIDUsed(ctx, "src", []string{})).To(Succeed())
		})
	})
})
