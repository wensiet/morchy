package usecase

import (
	"context"

	"github.com/wernsiet/morchy/pkg/agent/domain"
	"github.com/wernsiet/morchy/pkg/agent/domain/workload"
	apitypes "github.com/wernsiet/morchy/pkg/controlplane/implementation/jsonformatter"
	"go.uber.org/zap"
)

func (i *interactor) getRunnableWorkload(ctx context.Context) (*apitypes.WorkloadResponse, error) {
	availableWorkloads, err := i.controlplaneClient.ListAvailableWorkloads(ctx, *i.workloadRepo.GetResourceLimits())
	if err != nil {
		return nil, domain.ErrorBaseWorkloadInternal.Wrapf(err, "error getting available workloads from control-plane")
	}
	var chosenWorkload *apitypes.WorkloadResponse
	for _, candidate := range availableWorkloads {
		if storedWL, _ := i.workloadRepo.GetWorkload(candidate.ID); storedWL != nil {
			i.logger.Warn("skipped workload", zap.String(domain.SReason, domain.SWorkloadAlreadyInStoreage))
			continue
		}
		chosenWorkload = candidate
		break
	}
	return chosenWorkload, nil
}

func (i *interactor) CreateWorkload(ctx context.Context, chosenWorkload apitypes.WorkloadResponse) error {
	if err := i.controlplaneClient.CreateOrExtendWorkloadLease(ctx, chosenWorkload.ID); err != nil {
		return domain.ErrorBaseWorkloadInternal.With(domain.SWorkload, chosenWorkload.ID).
			Wrapf(err, "error leasing workload from control-plane")
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

	domainWorkload, err := i.workloadRepo.SaveWorkload(
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

	err = i.startWorklodLifecycle(ctx, *domainWorkload)
	if err != nil {
		domain.ErrorBaseWorkloadInternal.With(domain.SWorkload, chosenWorkload.ID).
			Wrapf(err, "error on starting wokrlod lifecycle")
	}

	return nil
}

func (i *interactor) terminateWorkload(ctx context.Context, w workload.Workload) error {
	i.logger.Info("terminating workload", zap.String(domain.SWorkloadID, w.ID))
	_ = i.runtimeClient.StopContainer(ctx, w.Container.Name)
	_ = i.runtimeClient.RemoveContainer(ctx, w.Container.Name) // TODO: push event or panic
	i.workloadRepo.RemoveWorkload(w.ID)
	return nil
}

func (i *interactor) ReconcileWorkload(ctx context.Context, wl workload.Workload) error {
	status, err := i.runtimeClient.GetContainerStatus(ctx, wl.Container.Name)
	if err != nil {
		return err
	}
	if status != domain.SRunning {
		i.stopWorkloadLifecycle(ctx, wl)
		err = i.terminateWorkload(ctx, wl) // TODO: push event
		if err != nil {
			i.logger.Error("unable to terminate workload", zap.Error(err))
		}
		return domain.ErrorBaseWorkloadHealthcheckFailed.
			With(domain.SWorkload, wl.ID).
			Errorf("container status: %s", status)
	}
	return i.controlplaneClient.CreateOrExtendWorkloadLease(ctx, wl.ID)
}
