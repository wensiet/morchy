package usecase

import (
	"context"
	"time"

	"github.com/wernsiet/morchy/pkg/agent/domain"
	"github.com/wernsiet/morchy/pkg/agent/domain/workload"
	"go.uber.org/zap"
)

func (i *interactor) ApplyWorkloadJoin(ctx context.Context) error {
	logger := i.logger.With(zap.String(domain.SUsecase, domain.SApplyWorkloadJoin))

	availableWorkloads, err := i.controlplaneClient.ListAvailableWorkloads(ctx, *i.workloadRepo.GetResourceLimits())
	if err != nil {
		return domain.ErrorBaseWorkloadInternal.Wrapf(err, "error getting available workloads from control-plane")
	}
	if len(availableWorkloads) == 0 {
		logger.Info("skipped workload join", zap.String(domain.SReason, domain.SNotWorkloadsToSchedule))
		return nil
	}
	chosenWorkload := availableWorkloads[0]
	logger = logger.With(zap.String("workloadID", chosenWorkload.ID))

	logger.Info("trying to join workload")

	if err := i.controlplaneClient.CreateWorkloadLease(ctx, chosenWorkload.ID); err != nil {
		// return domain.ErrorBaseWorkloadInternal.With(domain.SWorkload, chosenWorkload.ID).
		// 	Wrapf(err, "error leasing workload from control-plane")
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

	_, err = i.workloadRepo.SaveWorklod(*chosenWorkload)
	if err != nil {
		return domain.ErrorBaseWorkloadInternal.With(domain.SWorkload, chosenWorkload.ID).
			Wrapf(err, "error while saving started workload")
	}

	// Start async reconciliation loop to keep extending the lease
	go i.startWorkloadAsyncReconciliation(ctx, chosenWorkload)

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
			err := i.reconcileWorkload(ctx, w.ID)
			if err != nil {
				i.logger.Error("failed to reconcile workload", zap.Error(err))
				if terminateErr := i.terminateWorkload(ctx, *w); terminateErr != nil {
					i.logger.Error("failed to termintae workload", zap.Error(err))
				}
				return
			}
		}
	}
}

func (i *interactor) reconcileWorkload(ctx context.Context, workloadID string) error {
	wl, err := i.workloadRepo.GetWorkload(workloadID)
	if err != nil {
		return err
	}
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
