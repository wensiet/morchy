package usecase

import (
	"context"
	"time"

	"github.com/wernsiet/morchy/pkg/agent/domain"
	"github.com/wernsiet/morchy/pkg/agent/domain/workload"
	apitypes "github.com/wernsiet/morchy/pkg/controlplane/implementation/jsonformatter"
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
		domainWorkload, err := i.workloadRepo.SaveWorklod(
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
		loadedWorkloads = append(loadedWorkloads, domainWorkload)
	}

	for _, wl := range loadedWorkloads {
		err = i.ReconcileWorkload(ctx, *wl)
		if err != nil {
			return domain.ErrorBaseWorkloadInternal.With(domain.SWorkload, wl.ID).
				Wrapf(err, "error on loaded workload reconciliation")
		}
		go i.startWorkloadAsyncReconciliation(ctx, wl)
	}

	logger.Info("successfully loaded current state", zap.Int("loaded_workloads", len(loadedWorkloads)))

	return nil
}

func (i *interactor) ApplyWorkloadJoin(ctx context.Context) error {
	logger := i.logger.With(zap.String(domain.SUsecase, domain.SApplyWorkloadJoin))

	availableWorkloads, err := i.controlplaneClient.ListAvailableWorkloads(ctx, *i.workloadRepo.GetResourceLimits())
	if err != nil {
		return domain.ErrorBaseWorkloadInternal.Wrapf(err, "error getting available workloads from control-plane")
	}
	var chosenWorkload *apitypes.WorkloadResponse
	for _, candidate := range availableWorkloads {
		if storedWL, _ := i.workloadRepo.GetWorkload(candidate.ID); storedWL != nil {
			logger.Warn("skipped workload", zap.String(domain.SReason, domain.SWorkloadAlreadyInStoreage))
			continue
		}
		chosenWorkload = candidate
		break
	}
	if chosenWorkload == nil {
		logger.Info("skipped workload join", zap.String(domain.SReason, domain.SNoWorkloadsToSchedule))
		return nil
	}

	logger = logger.With(zap.String("workloadID", chosenWorkload.ID))

	logger.Info("trying to join workload")

	if err := i.controlplaneClient.CreateWorkloadLease(ctx, chosenWorkload.ID); err != nil {
		// return domain.ErrorBaseWorkloadInternal.With(domain.SWorkload, chosenWorkload.ID).
		// 	Wrapf(err, "error leasing workload from control-plane")
	}

	chosenWorkload.Container.Labels = map[string]string{
		domain.SWorkloadID: chosenWorkload.ID,
		domain.SManager:    domain.SAppName,
	}

	containerID, err := i.runtimeClient.CreateContainer(ctx, chosenWorkload.Container)
	if err != nil {
		return domain.ErrorBaseWorkloadInternal.With(domain.SWorkload, chosenWorkload.ID).
			Wrapf(err, "error while creating runtime container")
	}

	if err := i.runtimeClient.StartContainer(ctx, containerID); err != nil {
		return domain.ErrorBaseWorkloadInternal.With(domain.SWorkload, chosenWorkload.ID).
			Wrapf(err, "error while starting runtime container")
	}

	domainWorkload, err := i.workloadRepo.SaveWorklod(
		workload.Workload{
			ID: chosenWorkload.ID,
			Container: workload.Container{
				Name: chosenWorkload.Container.Name,
			},
		},
	)
	if err != nil {
		return domain.ErrorBaseWorkloadInternal.With(domain.SWorkload, chosenWorkload.ID).
			Wrapf(err, "error while saving started workload")
	}

	err = i.ReconcileWorkload(ctx, *domainWorkload)
	if err != nil {
		return domain.ErrorBaseWorkloadInternal.With(domain.SWorkload, chosenWorkload.ID).
			Wrapf(err, "error on initial workload reconciliation")
	}
	// Start async reconciliation loop to keep extending the lease
	go i.startWorkloadAsyncReconciliation(ctx, domainWorkload)

	logger.Info("successfully joined workload")

	return nil
}

func (i *interactor) terminateWorkload(ctx context.Context, w workload.Workload) error {
	i.logger.Info("terminating workload", zap.String("workloadID", w.ID))
	if err := i.runtimeClient.RemoveContainer(ctx, w.Container.Name); err != nil {
		return domain.ErrorBaseWorkloadInternal.With(domain.SWorkload, w.ID).
			Wrapf(err, "failed to remove container")
	}
	if err := i.runtimeClient.StopContainer(ctx, w.Container.Name); err != nil {
		return domain.ErrorBaseWorkloadInternal.With(domain.SWorkload, w.ID).
			Wrapf(err, "failed to stop container")
	}
	i.workloadRepo.RemoveWorkload(w.ID)
	return nil
}

func (i *interactor) startWorkloadAsyncReconciliation(ctx context.Context, w *workload.Workload) {
	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			_ = i.ReconcileWorkload(ctx, *w)
		}
	}
}

func (i *interactor) reconcileWorkload(ctx context.Context, wl workload.Workload) error {
	containerStatus, err := i.runtimeClient.GetContainerStatus(ctx, wl.Container.Name)
	if err != nil {
		return err
	}
	if containerStatus != domain.SRunning {
		return domain.ErrorBaseWorkloadHealthcheckFailed.With(domain.SWorkload, wl.ID).Errorf("workload healthcheck failed: %s", containerStatus)
	}
	if err = i.controlplaneClient.ExtendWorkloadLease(ctx, wl.ID); err != nil {
		return err
	}
	return nil
}

func (i *interactor) ReconcileWorkload(ctx context.Context, wl workload.Workload) error {
	if err := i.reconcileWorkload(ctx, wl); err != nil {
		i.logger.Error("failed to reconcile workload", zap.Error(err))
		if terminateErr := i.terminateWorkload(ctx, wl); terminateErr != nil {
			i.logger.Error("failed to termintae workload", zap.Error(err))
		}
		return err
	}
	return nil
}
