package usecase

import (
	"context"
	"time"

	"github.com/samber/oops"
	"github.com/wernsiet/morchy/pkg/agent/domain"
	"github.com/wernsiet/morchy/pkg/agent/domain/workload"
	apitypes "github.com/wernsiet/morchy/pkg/controlplane/implementation/jsonformatter"
	"github.com/wernsiet/morchy/pkg/runtime"
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

	err := i.runWorkload(ctx, chosenWorkload.Container)
	if err != nil {
		return domain.ErrorBaseWorkloadInternal.
			With(domain.SWorkload, chosenWorkload.ID).
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
	_ = i.runtimeClient.StopContainer(ctx, w.Container.Name)
	err := i.runtimeClient.RemoveContainer(ctx, w.Container.Name) // TODO: push event or panic
	if err != nil {
		i.logger.Error("failed to terminate workload", zap.String(domain.SWorkloadID, w.ID), zap.Error(err))
	} else {
		i.logger.Info("terminated workload", zap.String(domain.SWorkloadID, w.ID))
	}
	i.workloadRepo.RemoveWorkload(w.ID)
	return nil
}

func (i *interactor) runWorkload(ctx context.Context, workloadContainer runtime.Container) error {
	containerID, err := i.runtimeClient.CreateContainer(ctx, workloadContainer)
	if err != nil {
		return err
	}
	return i.runtimeClient.StartContainer(ctx, containerID)
}

func (i *interactor) ReconcileWorkload(ctx context.Context, wl workload.Workload) error {
	startTime := time.Now()
	status, err := i.runtimeClient.GetContainerStatus(ctx, wl.Container.Name)
	if err != nil {
		return err
	}

	err = func() error {
		if status != domain.SRunning {
			// TODO: push event
			return domain.ErrorBaseWorkloadHealthcheckFailed.
				With(domain.SWorkload, wl.ID).
				Errorf("container status: %s", status)
		}
		return i.controlplaneClient.CreateOrExtendWorkloadLease(ctx, wl.ID)
	}()

	if err != nil {
		oopsErr, ok := oops.AsOops(err)
		if !ok ||
			oopsErr.Code() == domain.SHealthcheckFailed ||
			oopsErr.Code() == domain.SOwnedByAnotherNode ||
			oopsErr.Code() == domain.STerminatedOnControlPlane {
			_ = i.stopWorkloadLifecycle(ctx, wl) // TODO: push event
		}
		return err
	}

	i.logger.Info(
		"reconciled workload",
		zap.String(domain.SWorkloadID, wl.ID),
		zap.String(domain.SWorkloadStatus, status),
		zap.Duration(domain.SDuration, time.Since(startTime)),
	)
	return nil
}
