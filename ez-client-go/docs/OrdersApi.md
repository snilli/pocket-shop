# {{classname}}

All URIs are relative to *https://api.ezcards.io*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateBulkOrder**](OrdersApi.md#CreateBulkOrder) | **Post** /v2/orders | Create Bulk Order
[**CreateInstantOrder**](OrdersApi.md#CreateInstantOrder) | **Post** /v2/orders/instant | Create Instant Order
[**GetOrders**](OrdersApi.md#GetOrders) | **Get** /v2/orders | Get Orders

# **CreateBulkOrder**
> InlineResponse2003 CreateBulkOrder(ctx, body)
Create Bulk Order

Create a new bulk order by specifying products and quantities. Operates asynchronously; fulfillment occurs subsequently. One API order can only contain products in a single currency. 

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**CreateBulkOrderRequest**](CreateBulkOrderRequest.md)|  | 

### Return type

[**InlineResponse2003**](inline_response_200_3.md)

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **CreateInstantOrder**
> InlineResponse2003 CreateInstantOrder(ctx, body)
Create Instant Order

Create a new order with exactly one product and quantity of one per request. The API will immediately attempt fulfillment. Returns COMPLETED if in stock (codes available via Get Codes), or PROCESSING if fulfilled asynchronously (poll order status). 

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**CreateInstantOrderRequest**](CreateInstantOrderRequest.md)|  | 

### Return type

[**InlineResponse2003**](inline_response_200_3.md)

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetOrders**
> InlineResponse2002 GetOrders(ctx, optional)
Get Orders

Retrieve an order list with pagination support. Use it to check the latest status of specific orders.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***OrdersApiGetOrdersOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a OrdersApiGetOrdersOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **limit** | **optional.Int32**|  | [default to 10]
 **page** | **optional.Int32**|  | [default to 1]
 **transactionId** | **optional.String**| Optional filter by transaction id for order lookups. | 
 **clientOrderNumber** | **optional.String**| Optional filter by clientOrderNumber for order lookups. | 

### Return type

[**InlineResponse2002**](inline_response_200_2.md)

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

