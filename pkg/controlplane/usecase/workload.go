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

func (i *interactor) identifyWorkloadStatusFromEvents(events []*workload.Event, lease *workload.Lease, stuckTimeout int) workload.WorkloadStatus {
	sort.Slice(events, func(i, j int) bool {
		return events[i].ProducedAt.Before(events[j].ProducedAt)
	})

	if lease == nil {
		if len(events) > 0 && events[len(events)-1].ProducedAt.Add(time.Duration(stuckTimeout)*time.Second).Before(time.Now()) {
			return workload.StuckWorkloadStatus
		}
		return workload.PendingWorkloadStatus
	}

	if len(events) < 2 {
		return workload.PendingWorkloadStatus
	}

	type briefHealthcheckPayload struct {
		Status string `json:"status"`
	}

	lastTwoEvents := events[len(events)-2:]
	successCount := 0
	failureCount := 0

	for _, event := range lastTwoEvents {
		var payload briefHealthcheckPayload
		_ = json.Unmarshal(event.Payload, &payload)
		if payload.Status == domain.SSuccess {
			successCount++
		} else {
			failureCount++
		}
	}

	if successCount == 2 {
		return workload.ActiveWorkloadStatus
	}

	if failureCount == 2 {
		return workload.FailedWorkloadStatus
	}

	return workload.DegradedWorkloadStatus
}

func (i *interactor) CreateWorkload(ctx context.Context, workloadSpec workload.WorkloadSpec) (*workload.Workload, error) {
	logger := i.logger.With(
		zap.String(domain.SDomain, domain.SWorkload),
		zap.String(domain.SContainerName, workloadSpec.Name),
	)

	workload, err := i.workloadRepo.CreateWorkload(ctx, workload.Workload{
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
	if err != nil {
		logger.Error("failed to start transaction", zap.Error(err))
		return nil, domain.ErrorWorkloadRepositoryInternalError.Wrap(err)
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != context.Canceled {
			logger.Debug("failed to rollback transaction", zap.Error(err))
		}
	}()

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
		i.eventListLimit,
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

	workload.Status = i.identifyWorkloadStatusFromEvents(events, lease, i.stuckTimeout)

	return workload, nil
}

func (i *interactor) ListWorkloads(ctx context.Context, statusEq *string, resourceLte *runtime.ResourceLimits, schedulableOnly bool) ([]*workload.Workload, error) {
	logger := i.logger.With(
		zap.String(domain.SDomain, domain.SWorkload),
	)

	workloads, err := i.workloadRepo.ListWorkloads(ctx, nil, resourceLte, schedulableOnly)
	if err != nil {
		logger.Error("failed to list workloads", zap.Error(err))
		return nil, err
	}

	for _, wl := range workloads {
		events, err := i.workloadRepo.ListEvents(
			ctx,
			map[string]string{
				domain.SAction:     domain.SHealthcheck,
				domain.SWorkloadID: wl.ID,
			},
			i.eventListLimit,
		)
		if err != nil {
			logger.Warn("failed to get workload events", zap.String(domain.SWorkloadID, wl.ID), zap.Error(err))
			continue
		}

		lease, err := i.workloadRepo.GetLeaseByWorkloadID(ctx, wl.ID)
		if err != nil {
			oopsErr, _ := oops.AsOops(err)
			if oopsErr.Code() != string(domain.NotFound) {
				logger.Warn("failed to get workload lease", zap.String(domain.SWorkloadID, wl.ID), zap.Error(err))
			}
		}

		if len(events) == 0 && lease == nil {
			continue
		}

		wl.Status = i.identifyWorkloadStatusFromEvents(events, lease, i.stuckTimeout)
	}

	if statusEq != nil {
		filtered := make([]*workload.Workload, 0)
		for _, wl := range workloads {
			if string(wl.Status) == *statusEq {
				filtered = append(filtered, wl)
			}
		}
		return filtered, nil
	}

	return workloads, nil
}

func (i *interactor) DeleteWorkload(ctx context.Context, workloadID string) error {
	logger := i.logger.With(
		zap.String(domain.SDomain, domain.SWorkload),
		zap.String(domain.SWorkloadID, workloadID),
	)

	if err := i.workloadRepo.DeleteWorkload(ctx, workloadID); err != nil {
		if oopsErr, ok := oops.AsOops(err); ok && oopsErr.Code() == string(domain.NotFound) {
			return err
		}
		logger.Error("failed to delete workload", zap.Error(err))
		return err
	}
	return nil
}
