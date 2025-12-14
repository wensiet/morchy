package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/samber/oops"
	"github.com/wernsiet/morchy/pkg/controlplane/domain"
	"github.com/wernsiet/morchy/pkg/controlplane/domain/workload"
	"github.com/wernsiet/morchy/pkg/runtime"
	"go.uber.org/zap"
)

func (i *interactor) CreateWorkload(ctx context.Context, workloadSpec workload.WorkloadSpec) (*workload.Workload, error) {
	logger := i.logger.With(
		zap.String(domain.SDomain, domain.SWorkload),
		zap.String(domain.SContainerName, workloadSpec.Name),
	)

	workload, err := i.wokrloadRepo.CreateWorkload(ctx, workload.Workload{
		ID:     uuid.NewString(),
		Status: workload.NewWorkloadStatus,
		Spec:   workloadSpec,
	})
	if err != nil {
		logger.Error("failed to create workload", zap.Error(err))
		return nil, err
	}
	return workload, nil
}

func (i *interactor) GetWorkload(ctx context.Context, workloadID string) (*workload.Workload, error) {
	logger := i.logger.With(
		zap.String(domain.SDomain, domain.SWorkload),
		zap.String(domain.SWorkloadID, workloadID),
	)

	workload, err := i.wokrloadRepo.GetWorkload(ctx, workloadID)
	if err != nil {
		if oopsErr, ok := oops.AsOops(err); ok && oopsErr.Code() == string(domain.NotFound) {
			return nil, err
		}
		logger.Error("failed to get workload", zap.Error(err))
		return nil, err
	}
	return workload, nil
}

func (i *interactor) ListWorkloads(ctx context.Context, statusEq *string, resourceLte *runtime.ResourceLimits) ([]*workload.Workload, error) {
	logger := i.logger.With(
		zap.String(domain.SDomain, domain.SWorkload),
	)

	workloads, err := i.wokrloadRepo.ListWorkloads(ctx, statusEq, resourceLte)
	if err != nil {
		logger.Error("failed to list workloads", zap.Error(err))
		return nil, err
	}
	return workloads, err
}
