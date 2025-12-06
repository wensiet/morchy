package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/wernsiet/morchy/pkg/agent/domain/workload"
	"github.com/wernsiet/morchy/pkg/runtime"
)

func (i *interactor) discoverWorkloads(ctx context.Context, limits runtime.ResourceLimits) ([]*workload.Workload, error) {
	type workloadResponse struct {
		ID     string `json:"id" example:"some-uuid"`
		Status string `json:"status" example:"new"`
	}
	u, err := url.Parse(i.controlPlaneURL + "/api/v1/workloads")
	if err != nil {
		return nil, fmt.Errorf("invalid url: %w", err)
	}

	q := u.Query()
	q.Set("cpu", fmt.Sprintf("%d", limits.CPU))
	q.Set("ram", fmt.Sprintf("%d", limits.RAM))
	q.Set("status", "new")
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := i.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	var workloads []workloadResponse
	if err := json.Unmarshal(body, &workloads); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	var result []*workload.Workload
	for _, wr := range workloads {
		result = append(result, &workload.Workload{
			ID:     wr.ID,
			Status: workload.WorkloadStatus(wr.Status),
		})
	}

	return result, nil
}

func (i *interactor) launchWorkload(ctx context.Context, w *workload.Workload) error {
	return nil
}

func (i *interactor) leaseWorkload(ctx context.Context, w *workload.Workload) error {
	return i.leaseMutationAction(ctx, w, http.MethodPost)
}

func (i *interactor) extendWorkloadLease(ctx context.Context, w *workload.Workload) error {
	return i.leaseMutationAction(ctx, w, http.MethodPut)
}

func (i *interactor) leaseMutationAction(ctx context.Context, w *workload.Workload, method string) error {
	u, err := url.Parse(i.controlPlaneURL + "/api/v1/workloads/" + w.ID + "/lease")
	if err != nil {
		return fmt.Errorf("invalid url: %w", err)
	}

	q := u.Query()
	q.Set("node_id", "some-node-uuid")
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, method, u.String(), nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := i.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Println("Raw JSON response:", string(body))

		return fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	return nil
}

func (i *interactor) discoverLimits(ctx context.Context) (*runtime.ResourceLimits, error) {
	return &runtime.ResourceLimits{
		CPU: 500,
		RAM: 1024,
	}, nil
}

func (i *interactor) ApplyWorkloadJoin(ctx context.Context) error {
	limits, err := i.discoverLimits(ctx)
	if err != nil {
		return err
	}
	availableWorkloads, err := i.discoverWorkloads(ctx, *limits)
	if err != nil {
		return err
	}
	if len(availableWorkloads) == 0 {
		return nil
	}
	chosenWorkload := availableWorkloads[0]
	fmt.Printf("Chosen workload: %s\n", chosenWorkload.ID)

	err = i.leaseWorkload(ctx, chosenWorkload)
	if err != nil {
		return err
	}

	err = i.launchWorkload(ctx, chosenWorkload)
	if err != nil {
		return err
	}

	return nil
}

func (i *interactor) ReconcileWorkloads(ctx context.Context) error {
	wl := &workload.Workload{
		ID:     "2a95b423-7586-41af-a591-c2fd1c4eb4a2",
		Status: workload.NewWorkloadStatus,
	}
	err := i.extendWorkloadLease(ctx, wl)
	if err != nil {
		return err
	}
	return nil
}
