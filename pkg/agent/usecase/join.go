package usecase

import (
	"context"

	"github.com/samber/oops"
	"github.com/wernsiet/morchy/pkg/agent/domain"
	"github.com/wernsiet/morchy/pkg/agent/domain/workload"
	"github.com/wernsiet/morchy/pkg/runtime"
	"go.uber.org/zap"
)

func (i *interactor) LoadCurrentState(ctx context.Context) error {
	logger := i.logger.With(zap.String(domain.SUsecase, domain.SLoadCurrentState))

	currentContainers, err := i.runtimeClient.ListContainers(ctx, &runtime.ContainerFilters{
		Labels: map[string]string{
			domain.SManager: domain.SAppName,
		},
	})
	if err != nil {
		return domain.ErrorBaseWorkloadInternal.
			Wrapf(err, "error while listing runtime containers")
	}

	loadedWorkloads := make([]*workload.Workload, 0, len(currentContainers))
	for _, container := range currentContainers {
		workloadID, ok := container.Labels["workload_id"]
		if !ok {
			return domain.ErrorBaseWorkloadInternal.
				New("container has invalid labels: workload_id label is missing")
		}
		domainWorkload, err := i.workloadRepo.SaveWorkload(
			workload.Workload{
				ID: workloadID,
				Container: workload.Container{
					Name: container.Name,
				},
			},
		)
		if err != nil {
			return domain.ErrorBaseWorkloadInternal.
				Wrapf(err, "error while saving runtime containers")
		}

		if err := i.startWorklodLifecycle(ctx, *domainWorkload); err != nil {
			oopsErr, ok := oops.AsOops(err)
			if !ok || oopsErr.Code() != domain.SHealthcheckFailed || oopsErr.Code() != domain.STerminatedOnControlPlane {
				return domain.ErrorBaseWorkloadInternal.With(domain.SWorkload, domainWorkload.ID).
					Wrapf(err, "error on starting workload lifecycle")
			}
		}

		loadedWorkloads = append(loadedWorkloads, domainWorkload)
	}

	logger.Info("successfully loaded current state", zap.Int("loaded_workloads", len(loadedWorkloads)))

	return nil
}

func (i *interactor) ApplyWorkloadJoin(ctx context.Context) error {
	chosenWorkload, err := i.getRunnableWorkload(ctx)
	if err != nil {
		return domain.ErrorBaseWorkloadInternal.
			Wrapf(err, "failed to get runnable workload from controlplane")
	}
	if chosenWorkload == nil {
		return nil
	}

	err = i.CreateWorkload(ctx, *chosenWorkload)
	if err != nil {
		return err
	}

	i.logger.Info("successfully joined workload", zap.String(domain.SUsecase, domain.SApplyWorkloadJoin), zap.String(domain.SWorkloadID, chosenWorkload.ID))

	return nil
}
