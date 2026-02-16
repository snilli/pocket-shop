package ezclient

import (
	swagger "ez-client-go"

	"pocket-shop/internal/port"
)

func ToModel(transactionID, clientOrderNumber string, status *swagger.OrderStatus) *port.Order {
	st := swagger.PROCESSING_OrderStatus
	if status != nil && *status != "" {
		st = *status
	}
	return &port.Order{
		TransactionId:     transactionID,
		ClientOrderNumber: clientOrderNumber,
		Status:            port.OrderStatusDTO(st),
	}
}

func ToModelFromOrderResponse(order *swagger.OrderResponse, fallbackClientOrderNumber string) *port.Order {
	if order == nil {
		return ToModel("", fallbackClientOrderNumber, nil)
	}
	return ToModel(order.TransactionId, order.ClientOrderNumber, order.Status)
}

func ToModelFromGetOrdersResponse(data *swagger.OrdersPaginationResponse, fallbackClientOrderNumber string) *port.Order {
	if data == nil || len(data.Items) == 0 {
		return ToModel("", fallbackClientOrderNumber, nil)
	}
	return ToModelFromGetOrdersItem(&data.Items[0])
}

func ToModelFromGetOrdersItem(order *swagger.GetOrdersOrderResponse) *port.Order {
	if order == nil {
		return ToModel("", "", nil)
	}
	return ToModel(order.TransactionId, order.ClientOrderNumber, order.Status)
}
