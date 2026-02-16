package domain

import (
	"time"
)

type OrderStatus string

const (
	StatusProcessing OrderStatus = "PROCESSING"
	StatusCompleted  OrderStatus = "COMPLETED"
	StatusCancelled  OrderStatus = "CANCELLED"
)

type Order struct {
	ID        string      `json:"id"`
	RefID     string      `json:"-"`
	RefSource string      `json:"-"`
	Status    OrderStatus `json:"status"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt *time.Time  `json:"updatedAt,omitempty"`
}

func (o *Order) updatedAt() {
	now := time.Now()
	o.UpdatedAt = &now
}

func (o *Order) Complete() {
	o.Status = StatusCompleted
	o.updatedAt()
}

func (o *Order) Cancel() {
	o.Status = StatusCancelled
	o.updatedAt()
}

func Create(refID, refSource string) *Order {
	return &Order{
		RefID:     refID,
		RefSource: refSource,
		Status:    StatusProcessing,
		CreatedAt: time.Now(),
	}
}
