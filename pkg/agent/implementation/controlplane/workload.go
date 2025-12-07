package controlplane

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/wernsiet/morchy/pkg/agent/domain"
	"github.com/wernsiet/morchy/pkg/agent/domain/workload"
	"github.com/wernsiet/morchy/pkg/runtime"
)

func (c *Client) ListAvailableWorkloads(ctx context.Context, limits runtime.ResourceLimits) ([]*workload.Workload, error) {
	type workloadResponse struct {
		ID        string            `json:"id"`
		Status    string            `json:"status"`
		Container runtime.Container `json:"container"`
	}

	var workloads []workloadResponse
	request, err := c.httpClient.R().
		SetContext(ctx).
		SetHeader("Accept", "application/json").
		SetQueryParam("cpu", fmt.Sprintf("%d", limits.CPU)).
		SetQueryParam("ram", fmt.Sprintf("%d", limits.RAM)).
		SetQueryParam("status", "new").
		SetResult(&workloads).
		Get(c.baseURL + "/api/v1/workloads")
	if err != nil {
		return nil, err
	}
	if request.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", request.StatusCode())
	}

	result := make([]*workload.Workload, 0, len(workloads))
	for _, wr := range workloads {
		result = append(result, &workload.Workload{
			ID:        wr.ID,
			Status:    workload.WorkloadStatus(wr.Status),
			Container: wr.Container,
		})
	}

	return result, nil
}

func (c *Client) leaseMutationAction(ctx context.Context, workloadID, method string) error {
	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetHeader("Accept", "application/json").
		SetQueryParam("node_id", "some-node-uuid"). // TODO: replace set node id on start and store in interactor
		Execute(method, c.baseURL+"/api/v1/workloads/"+workloadID+"/lease")
	if err != nil {
		return domain.ErrorBaseWorkloadInternal.Join(err)
	}

	if statusCode := resp.StatusCode(); statusCode != http.StatusOK {
		if statusCode == http.StatusNotFound {
			return domain.ErrorWorkloadTerminatedOnControlPlane.Errorf("workload %s was not found on control-plane, treat it as terminated", workloadID)
		}
		if statusCode == http.StatusConflict && strings.Contains(resp.String(), domain.SOwnedByAnotherNode) {
			return domain.ErrorWorkloadOwnedByAnotherNode.Errorf("workload %s is owned by another node", workloadID)
		}
		return domain.ErrorBaseWorkloadInternal.Errorf("got unexpected status code for workload id=%s lease extension", workloadID)
	}

	return nil
}

func (c *Client) CreateWorkloadLease(ctx context.Context, workloadID string) error {
	return c.leaseMutationAction(ctx, workloadID, http.MethodPost)
}

func (c *Client) ExtendWorkloadLease(ctx context.Context, workloadID string) error {
	return c.leaseMutationAction(ctx, workloadID, http.MethodPut)
}
