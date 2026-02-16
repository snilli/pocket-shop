package port

type OrderStatusDTO string

const (
	PROCESSING_OrderStatus OrderStatusDTO = "PROCESSING"
	COMPLETED_OrderStatus  OrderStatusDTO = "COMPLETED"
	CANCELLED_OrderStatus  OrderStatusDTO = "CANCELLED"
)

type Order struct {
	TransactionId     string
	ClientOrderNumber string
	Status            OrderStatusDTO
}

func (o *Order) IsProcessing() bool {
	return o.Status == PROCESSING_OrderStatus
}

func (o *Order) IsCompleted() bool {
	return o.Status == COMPLETED_OrderStatus
}

func (o *Order) IsCancelled() bool {
	return o.Status == CANCELLED_OrderStatus
}

func (o *Order) GetStatus() string {
	return string(o.Status)
}
