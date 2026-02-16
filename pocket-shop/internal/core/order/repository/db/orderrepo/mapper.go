package orderrepo

import (
	"orm/ent"
	"time"

	entorder "orm/ent/order"

	"pocket-shop/internal/core/order/domain"
)

func toDomain(e *ent.Order) *domain.Order {
	if e == nil {
		return nil
	}
	return &domain.Order{
		ID:        e.ID.String(),
		RefID:     e.RefID,
		RefSource: e.RefSource,
		Status:    fromEntStatus(e.Status),
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

type orderModel struct {
	ID        string
	RefID     string
	RefSource string
	Status    entorder.Status
	CreatedAt time.Time
	UpdatedAt *time.Time
}

func toModel(o domain.Order) orderModel {
	return orderModel{
		ID:        o.ID,
		RefID:     o.RefID,
		RefSource: o.RefSource,
		Status:    toEntStatus(o.Status),
		CreatedAt: o.CreatedAt,
		UpdatedAt: o.UpdatedAt,
	}
}

func toEntStatus(s domain.OrderStatus) entorder.Status {
	switch s {
	case domain.StatusProcessing:
		return entorder.StatusPROCESSING
	case domain.StatusCompleted:
		return entorder.StatusCOMPLETED
	case domain.StatusCancelled:
		return entorder.StatusCANCELLED
	default:
		return entorder.StatusPROCESSING
	}
}

func fromEntStatus(s entorder.Status) domain.OrderStatus {
	switch s {
	case entorder.StatusPROCESSING:
		return domain.StatusProcessing
	case entorder.StatusCOMPLETED:
		return domain.StatusCompleted
	case entorder.StatusCANCELLED:
		return domain.StatusCancelled
	default:
		return domain.StatusProcessing
	}
}
