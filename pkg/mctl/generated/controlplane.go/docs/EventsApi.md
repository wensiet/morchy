# {{classname}}

All URIs are relative to */*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ApiV1EventsPost**](EventsApi.md#ApiV1EventsPost) | **Post** /api/v1/events | Push event

# **ApiV1EventsPost**
> ApiV1EventsPost(ctx, body, nodeId)
Push event

Create a new event for a node

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**JsonformatterEventCreateRequest**](JsonformatterEventCreateRequest.md)| Event payload | 
  **nodeId** | **string**| Node ID | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

