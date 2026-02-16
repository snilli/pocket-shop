# {{classname}}

All URIs are relative to *https://api.ezcards.io*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetExchangeRates**](ExchangeRatesApi.md#GetExchangeRates) | **Get** /v2/exchange-rates | Get Exchange Rates

# **GetExchangeRates**
> InlineResponse2005 GetExchangeRates(ctx, optional)
Get Exchange Rates

Provide access to the latest EZ Rate (effective exchange rate) for supported currency pairs. Only the EZ Rate is returned; market rates and internal buffers are not exposed.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***ExchangeRatesApiGetExchangeRatesOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a ExchangeRatesApiGetExchangeRatesOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **pairs** | **optional.String**| Comma-separated list of currency pairs. Each pair is 6 letters in BASEQUOTE ISO-4217 format (e.g., USDCAD,USDEUR,GBPEUR). If provided, base and quote are ignored. | 
 **base** | **optional.String**| Filter by base currency (ISO-4217), e.g. USD. Ignored if pairs is provided. | 
 **quote** | **optional.String**| Filter by quote currency (ISO-4217), e.g. CAD. Ignored if pairs is provided. | 

### Return type

[**InlineResponse2005**](inline_response_200_5.md)

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

