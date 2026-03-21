# {{classname}}

All URIs are relative to */*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ApiV1WorkloadsWorkloadIdLeaseDelete**](LeasesApi.md#ApiV1WorkloadsWorkloadIdLeaseDelete) | **Delete** /api/v1/workloads/{workload_id}/lease | Release a lease
[**ApiV1WorkloadsWorkloadIdLeaseGet**](LeasesApi.md#ApiV1WorkloadsWorkloadIdLeaseGet) | **Get** /api/v1/workloads/{workload_id}/lease | Get lease for workload
[**ApiV1WorkloadsWorkloadIdLeasePut**](LeasesApi.md#ApiV1WorkloadsWorkloadIdLeasePut) | **Put** /api/v1/workloads/{workload_id}/lease | Extend a lease

# **ApiV1WorkloadsWorkloadIdLeaseDelete**
> ApiV1WorkloadsWorkloadIdLeaseDelete(ctx, workloadId, nodeId)
Release a lease

Release/delete a lease for a specific workload on a given node

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **workloadId** | **string**| Workload ID | 
  **nodeId** | **string**| Node ID | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ApiV1WorkloadsWorkloadIdLeaseGet**
> JsonformatterLeaseResponse ApiV1WorkloadsWorkloadIdLeaseGet(ctx, workloadId)
Get lease for workload

Get the current lease for a specific workload

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **workloadId** | **string**| Workload ID | 

### Return type

[**JsonformatterLeaseResponse**](jsonformatter.LeaseResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ApiV1WorkloadsWorkloadIdLeasePut**
> JsonformatterLeaseResponse ApiV1WorkloadsWorkloadIdLeasePut(ctx, workloadId, nodeId)
Extend a lease

Extend a lease for a specific workload on a given node

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **workloadId** | **string**| Workload ID | 
  **nodeId** | **string**| Node ID | 

### Return type

[**JsonformatterLeaseResponse**](jsonformatter.LeaseResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

