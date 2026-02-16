# {{classname}}

All URIs are relative to *https://api.ezcards.io*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetCodes**](CodesApi.md#GetCodes) | **Get** /v2/orders/{transactionId}/codes | Get Codes

# **GetCodes**
> InlineResponse2004 GetCodes(ctx, transactionId, optional)
Get Codes

Fetch codes that are ready for use within specific orders. Codes become available progressively as the order is fulfilled.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **transactionId** | **string**| Unique identifier of the transaction to retrieve codes for. | 
 **optional** | ***CodesApiGetCodesOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a CodesApiGetCodesOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **sku** | **optional.String**| Optional filter by product SKU. | 
 **stockIds** | **optional.String**| Optional filter by list of EZ stock identifiers. | 

### Return type

[**InlineResponse2004**](inline_response_200_4.md)

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

