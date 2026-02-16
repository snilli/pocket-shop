# GetOrdersOrderResponse

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**TransactionId** | **string** | The transaction identifier. | [default to null]
**ClientOrderNumber** | **string** | The client&#x27;s order number. | [default to null]
**EzOrderNumber** | **string** | EZ order number. Appears in doc example only; not listed in schema table. | [optional] [default to null]
**GrandTotal** | [***MoneyResponse**](MoneyResponse.md) |  | [default to null]
**CreatedAt** | [**time.Time**](time.Time.md) | Timestamp of when the order was created. | [default to null]
**Status** | [***OrderStatus**](OrderStatus.md) |  | [default to null]
**Products** | [**[]OrderProductResponseGetOrders**](OrderProductResponseGetOrders.md) | List of products in the order. | [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

