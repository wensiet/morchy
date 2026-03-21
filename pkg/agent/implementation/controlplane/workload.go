package controlplane

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/wernsiet/morchy/pkg/agent/domain"
	apitypes "github.com/wernsiet/morchy/pkg/controlplane/implementation/jsonformatter"
	"github.com/wernsiet/morchy/pkg/runtime"
)

func (c *Client) ListAvailableWorkloads(ctx context.Context, limits runtime.ResourceLimits) ([]*apitypes.WorkloadResponse, error) {
	var workloads []*apitypes.WorkloadResponse
	request, err := c.httpClient.R().
		SetContext(ctx).
		SetHeader("Accept", "application/json").
		SetQueryParam("cpu", fmt.Sprintf("%d", limits.CPU)).
		SetQueryParam("ram", fmt.Sprintf("%d", limits.RAM)).
		SetQueryParam("schedulable_only", "true").
		SetResult(&workloads).
		Get(c.baseURL + "/api/v1/workloads")
	if err != nil {
		return nil, err
	}
	if request.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", request.StatusCode())
	}

	return workloads, nil
}

func (c *Client) CreateOrExtendWorkloadLease(ctx context.Context, workloadID string) error {
	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetHeader("Accept", "application/json").
		SetQueryParam("node_id", c.nodeID).
		Execute(http.MethodPut, c.baseURL+"/api/v1/workloads/"+workloadID+"/lease")
	if err != nil {
		return domain.ErrorBaseWorkloadInternal.Wrap(err)
	}

	if statusCode := resp.StatusCode(); statusCode != http.StatusOK {
		if statusCode == http.StatusNotFound {
			return domain.ErrorWorkloadTerminatedOnControlPlane.Errorf("workload %s was not found on control-plane, treat it as terminated", workloadID)
		}
		if statusCode == http.StatusInternalServerError && strings.Contains(resp.String(), "lease_workload_id_fkey") {
			return domain.ErrorWorkloadTerminatedOnControlPlane.Errorf("workload %s was not found on control-plane, treat it as terminated", workloadID)
		}
		if statusCode == http.StatusConflict && strings.Contains(resp.String(), domain.SOwnedByAnotherNode) {
			return domain.ErrorWorkloadOwnedByAnotherNode.Errorf("workload %s is owned by another node", workloadID)
		}
		return domain.ErrorBaseWorkloadInternal.Errorf("got unexpected status code for workload id=%s lease extension", workloadID)
	}

	return nil
}

func (c *Client) DeleteWorkloadLease(ctx context.Context, workloadID string) error {
	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetHeader("Accept", "application/json").
		SetQueryParam("node_id", c.nodeID).
		Execute(http.MethodDelete, c.baseURL+"/api/v1/workloads/"+workloadID+"/lease")
	if err != nil {
		return domain.ErrorBaseWorkloadInternal.Wrap(err)
	}

	if statusCode := resp.StatusCode(); statusCode != http.StatusNoContent {
		if statusCode == http.StatusNotFound {
			return domain.ErrorWorkloadTerminatedOnControlPlane.Errorf("workload %s was not found on control-plane, treat it as terminated", workloadID)
		}
		if statusCode == http.StatusConflict && strings.Contains(resp.String(), domain.SOwnedByAnotherNode) {
			return domain.ErrorWorkloadOwnedByAnotherNode.Errorf("workload %s is owned by another node", workloadID)
		}
		return domain.ErrorBaseWorkloadInternal.Errorf("got unexpected status code for workload id=%s lease deletion", workloadID)
	}

	return nil
}

func (c *Client) PushEvent(ctx context.Context, event apitypes.EventCreateRequest) error {
	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetHeader("Accept", "application/json").
		SetQueryParam("node_id", c.nodeID).
		SetBody(event).
		Execute(http.MethodPost, c.baseURL+"/api/v1/events")
	if err != nil {
		return domain.ErrorBaseWorkloadInternal.Wrap(err)
	}

	if resp.StatusCode() != http.StatusCreated {
		return domain.ErrorBaseWorkloadInternal.Errorf("unexpected status: %d", resp.StatusCode())
	}

	return nil
}
