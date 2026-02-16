# ExchangeRatesPayload

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**UpdatedAt** | [**time.Time**](time.Time.md) | ISO-8601 timestamp of the latest rate update (UTC). | [default to null]
**ProductCurrencies** | **[]string** | Platform-supported currencies (ISO-4217). | [default to null]
**Rates** | [**[]ExchangeRate**](ExchangeRate.md) | List of EZ Rates for currency pairs. | [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

