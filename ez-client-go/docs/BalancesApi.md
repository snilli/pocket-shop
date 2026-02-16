# {{classname}}

All URIs are relative to *https://api.ezcards.io*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetBalances**](BalancesApi.md#GetBalances) | **Get** /v2/balances | Get Balances

# **GetBalances**
> InlineResponse200 GetBalances(ctx, optional)
Get Balances

Provide access to balance information across various currencies.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***BalancesApiGetBalancesOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a BalancesApiGetBalancesOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **currency** | **optional.String**| Optional filter for balance information by currency. | 

### Return type

[**InlineResponse200**](inline_response_200.md)

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

