# {{classname}}

All URIs are relative to *https://api.ezcards.io*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetProducts**](ProductsApi.md#GetProducts) | **Get** /v2/products | Get Products

# **GetProducts**
> InlineResponse2001 GetProducts(ctx, optional)
Get Products

Retrieve detailed listings of products, including specifications, pricing, availability, and other critical information.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***ProductsApiGetProductsOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a ProductsApiGetProductsOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **limit** | **optional.Int32**|  | [default to 1000]
 **page** | **optional.Int32**|  | [default to 1]
 **sku** | **optional.String**| Optional filter for listing products by their SKU (e.g. 8PX-UF-Y5U,TTR-WI-YWN). | 

### Return type

[**InlineResponse2001**](inline_response_200_1.md)

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

