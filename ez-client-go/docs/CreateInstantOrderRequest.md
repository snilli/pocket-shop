# CreateInstantOrderRequest

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ClientOrderNumber** | **string** | The order number assigned by the client, used for tracking and reference purposes. | [optional] [default to null]
**EnableClientOrderNumberDupCheck** | **bool** | Enable to prevent duplicate clientOrderNumber. | [optional] [default to false]
**Sku** | **string** | The Stock Keeping Unit identifier for the product. | [default to null]
**PayWithCurrency** | **string** | Settlement currency (ISO-4217). Required only when the product currency is not in the customer&#x27;s balance currencies.  | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

