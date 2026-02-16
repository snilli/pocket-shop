# ProductResponse

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Sku** | **string** | The Stock Keeping Unit identifier for the product. | [default to null]
**Name** | **string** | The name of the product. | [default to null]
**Brand** | **string** | The brand of the product. | [default to null]
**BrandCategory** | **string** | The category of the brand. | [optional] [default to null]
**BrandSubCategory** | **string** | The sub-category of the brand. | [optional] [default to null]
**Types** | **[]string** | List of product types. | [default to null]
**Format** | [***ProductFormat**](ProductFormat.md) |  | [default to null]
**Country** | **string** | The country for which the product is available. | [optional] [default to null]
**Currency** | **string** | The 3-letter currency code (ISO-4217). | [default to null]
**ImageUrl** | **string** | The URL of the image representing the product. | [optional] [default to null]
**Prices** | [**[]PriceResponse**](PriceResponse.md) |  | [default to null]
**FaceValue** | **string** | The face value of the product. | [default to null]
**PercentageOffFaceValue** | **string** | The percentage off the face value of the product. | [default to null]
**IsInstantDeliverySupported** | **bool** | Whether the CreateInstantOrder API can be used for instant order and fulfillment. | [default to null]
**Descriptions** | **[]string** |  | [default to null]
**Instructions** | **[]string** |  | [default to null]
**TermConditions** | **[]string** |  | [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

