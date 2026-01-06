package usecase

import (
	"context"
	"encoding/json"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/samber/oops"
	"github.com/wernsiet/morchy/pkg/controlplane/domain"
	"github.com/wernsiet/morchy/pkg/controlplane/domain/workload"
	"github.com/wernsiet/morchy/pkg/runtime"
	"go.uber.org/zap"
)

func identifyWorkloadStatusFromEvents(events []*workload.Event, lease *workload.Lease) workload.WorkloadStatus {
	/*
		Gather status strategy:
		- PENDING if workload is created, has no leases
		- STUCK if workload has no leases more than N seconds
		- ACTIVE if workload is running, with last max(N) success healthchecks
		- FAILED if workload is running, with last max(N) failed healthcheckes
		- DEGRADED (extra) if worklad is running, with flaping healtchecks
	*/
	sort.Slice(events, func(i, j int) bool {
		return events[i].ProducedAt.Before(events[j].ProducedAt)
	})

	if lease == nil {
		if len(events) > 0 && events[len(events)-1].ProducedAt.Add(time.Duration(domain.CStuckTimeout)*time.Second).Before(time.Now()) {
			return workload.StuckWorkloadStatus
		}
		return workload.PendingWorkloadStatus
	}

	lastStatus := workload.PendingWorkloadStatus

	type briefHealthcheckPayload struct {
		Status string `json:"status"`
	}

	totalFailedStatuses := 0
	for _, event := range events {
		var payload *briefHealthcheckPayload
		_ = json.Unmarshal(event.Payload, &payload)
		if payload.Status != domain.SSuccess {
			totalFailedStatuses += 1
		}
	}

	if totalFailedStatuses != 0 && totalFailedStatuses == len(events)/2 {
		return workload.DegradedWorkloadStatus
	}

	return lastStatus
}

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

	tx, err := i.dbPool.Begin(ctx)
	defer tx.Rollback(ctx)
	if err != nil {
		logger.Error("failed to start transaction", zap.Error(err))
		return nil, domain.ErrorWorkloadRepositoryInternalError.Wrap(err)
	}

	repo := i.repositoryFactory.New(tx)

	workload, err := repo.GetWorkload(ctx, workloadID)
	if err != nil {
		if oopsErr, ok := oops.AsOops(err); ok && oopsErr.Code() == string(domain.NotFound) {
			return nil, err
		}
		logger.Error("failed to get workload", zap.Error(err))
		return nil, err
	}
	events, err := repo.ListEvents(
		ctx,
		map[string]string{
			domain.SAction:     domain.SHealthcheck,
			domain.SWorkloadID: workloadID,
		},
		domain.CEventListLimit,
	)
	if err != nil {
		logger.Error("failed to get workload events", zap.Error(err))
		return nil, err
	}
	lease, err := repo.GetLeaseByWorkloadID(ctx, workloadID)
	if err != nil {
		oopsErr, _ := oops.AsOops(err)
		if oopsErr.Code() != string(domain.NotFound) {
			logger.Error("failed to get workload lease", zap.Error(err))
			return nil, err
		}
	}

	workload.Status = identifyWorkloadStatusFromEvents(events, lease)

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

func (i *interactor) DeleteWorkload(ctx context.Context, workloadID string) error {
	logger := i.logger.With(
		zap.String(domain.SDomain, domain.SWorkload),
		zap.String(domain.SWorkloadID, workloadID),
	)

	if err := i.wokrloadRepo.DeleteWorkload(ctx, workloadID); err != nil {
		if oopsErr, ok := oops.AsOops(err); ok && oopsErr.Code() == string(domain.NotFound) {
			return err
		}
		logger.Error("failed to delete workload", zap.Error(err))
		return err
	}
	return nil
}
