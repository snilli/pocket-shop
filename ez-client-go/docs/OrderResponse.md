# OrderResponse

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**TransactionId** | **string** | The transaction identifier. | [default to null]
**ClientOrderNumber** | **string** | The client&#x27;s order number. | [default to null]
**PayWithCurrency** | **string** | The currency used to settle/charge the order (Create Bulk/Instant only). | [optional] [default to null]
**GrandTotal** | [***MoneyResponse**](MoneyResponse.md) |  | [default to null]
**CreatedAt** | [**time.Time**](time.Time.md) | Timestamp of when the order was created. | [default to null]
**Status** | [***OrderStatus**](OrderStatus.md) |  | [default to null]
**Products** | [**[]OrderProductResponse**](OrderProductResponse.md) |  | [default to null]
**FxSummary** | [**[]FxEntry**](FxEntry.md) | Summary of FX pairs/rates when cross-currency conversion is applied. Omitted for same-currency orders. | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

