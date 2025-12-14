# {{classname}}

All URIs are relative to */*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ApiV1WorkloadsGet**](WorkloadsApi.md#ApiV1WorkloadsGet) | **Get** /api/v1/workloads | List workloads
[**ApiV1WorkloadsPost**](WorkloadsApi.md#ApiV1WorkloadsPost) | **Post** /api/v1/workloads | Create workload
[**ApiV1WorkloadsWorkloadIdDelete**](WorkloadsApi.md#ApiV1WorkloadsWorkloadIdDelete) | **Delete** /api/v1/workloads/{workload_id} | Delete workload
[**ApiV1WorkloadsWorkloadIdGet**](WorkloadsApi.md#ApiV1WorkloadsWorkloadIdGet) | **Get** /api/v1/workloads/{workload_id} | Get workload

# **ApiV1WorkloadsGet**
> []JsonformatterWorkloadResponse ApiV1WorkloadsGet(ctx, optional)
List workloads

Retrieve a list of workloads filtered by status, CPU, and RAM

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***WorkloadsApiApiV1WorkloadsGetOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a WorkloadsApiApiV1WorkloadsGetOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **status** | **optional.String**| Filter by workload status | 
 **cpu** | **optional.Int32**| Filter by CPU (millicores) | 
 **ram** | **optional.Int32**| Filter by RAM (MB) | 

### Return type

[**[]JsonformatterWorkloadResponse**](jsonformatter.WorkloadResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ApiV1WorkloadsPost**
> JsonformatterWorkloadResponse ApiV1WorkloadsPost(ctx, body)
Create workload

Create a new workload from provided spec

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**JsonformatterWorkloadSpecRequest**](JsonformatterWorkloadSpecRequest.md)| Workload specification | 

### Return type

[**JsonformatterWorkloadResponse**](jsonformatter.WorkloadResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ApiV1WorkloadsWorkloadIdDelete**
> ApiV1WorkloadsWorkloadIdDelete(ctx, workloadId)
Delete workload

Delete a workload and all related models (lease, spec)

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **workloadId** | **string**| Workload ID | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ApiV1WorkloadsWorkloadIdGet**
> JsonformatterWorkloadResponse ApiV1WorkloadsWorkloadIdGet(ctx, workloadId)
Get workload

Retrieve a single workload by ID

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **workloadId** | **string**| Workload ID | 

### Return type

[**JsonformatterWorkloadResponse**](jsonformatter.WorkloadResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

