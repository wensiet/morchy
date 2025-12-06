package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/wernsiet/morchy/internal/domain/workload"
	pkgworkload "github.com/wernsiet/morchy/pkg/workload"
)

func (i *interactor) CreateWorkload(ctx context.Context, workloadSpec workload.WorkloadSpec) (*workload.Workload, error) {
	workload, err := i.wokrloadRepo.CreateWorklod(ctx, workload.Workload{
		ID:     uuid.NewString(),
		Status: workload.NewWorkloadStatus,
		Spec:   workloadSpec,
	})
	if err != nil {
		return nil, err
	}
	return workload, nil
}

func (i *interactor) GetWorkload(ctx context.Context, workloadID string) (*workload.Workload, error) {
	workload, err := i.wokrloadRepo.GetWorkload(ctx, workloadID)
	if err != nil {
		return nil, err
	}
	return workload, nil
}

func (i *interactor) ListWorkloads(ctx context.Context, statusEq *string, resourceLte *pkgworkload.ResourceLimits) ([]*workload.Workload, error) {
	workloads, err := i.wokrloadRepo.ListWorkloads(ctx, statusEq, resourceLte)
	if err != nil {
		return nil, err
	}
	return workloads, err
}
