package usecase

import (
	"context"
	"encoding/json"
	"time"

	"github.com/samber/oops"
	"github.com/wernsiet/morchy/pkg/controlplane/domain"
	"github.com/wernsiet/morchy/pkg/controlplane/domain/workload"
	"go.uber.org/zap"
)

func (i *interactor) GetLeaseByNodeAndWorkloadID(ctx context.Context, nodeId, workloadId string) (*workload.Lease, error) {
	logger := i.logger.With(
		zap.String(domain.SDomain, domain.SWorkload),
		zap.String(domain.SNodeID, nodeId),
		zap.String(domain.SWorkloadID, workloadId),
	)

	lease, err := i.wokrloadRepo.GetLease(ctx, nodeId, workloadId)
	if err != nil {
		if oopsErr, ok := oops.AsOops(err); ok && oopsErr.Code() == string(domain.NotFound) {
			return nil, err
		}
		logger.Error("failed to get lease", zap.Error(err))
		return nil, err
	}
	return lease, nil
}

func (i *interactor) CreateOrExtendLease(ctx context.Context, nodeId, workloadId string) (*workload.Lease, error) {
	logger := i.logger.With(
		zap.String(domain.SDomain, domain.SWorkload),
		zap.String(domain.SNodeID, nodeId),
		zap.String(domain.SWorkloadID, workloadId),
	)

	tx, err := i.dbPool.Begin(ctx)
	defer tx.Rollback(ctx)
	if err != nil {
		logger.Error("failed to start transaction", zap.Error(err))
		return nil, domain.ErrorWorkloadRepositoryInternalError.Wrap(err)
	}

	repo := i.repositoryFactory.New(tx)

	lease, err := repo.UpsertLease(ctx, nodeId, workloadId)
	if err != nil {
		logger.Error("failed to upsert lease", zap.Error(err))
		return nil, err
	}

	// If lease created_at != updated_at it means that it was updated
	// and we can say that workload healthcheck is successful
	if !lease.CreatedAt.Truncate(time.Second).Equal(lease.UpdatedAt.Truncate(time.Second)) {
		eventPayload := map[string]string{
			domain.SWorkloadID: workloadId,
			domain.SAction:     domain.SHealthcheck,
			domain.SStatus:     domain.SSuccess,
		}
		jsonPayload, _ := json.Marshal(eventPayload)
		repo.SaveEvent(
			ctx, newEvent(nodeId, jsonPayload),
		)
	}

	err = tx.Commit(ctx)
	if err != nil {
		logger.Error("failed to commit transaction", zap.Error(err))
		return nil, domain.ErrorWorkloadRepositoryInternalError.Wrap(err)
	}

	logger.Info("upserted lease", zap.Time(domain.SUpdatedAt, lease.UpdatedAt))

	return lease, nil
}

func (i *interactor) ExpireLeases(ctx context.Context) error {
	logger := i.logger.With(
		zap.String(domain.SDomain, domain.SWorkload),
	)

	err := i.wokrloadRepo.DeleteExpiredLeases(ctx, domain.CLeaseLifetime)
	if err != nil {
		logger.Error("failed to delete expired leases", zap.Error(err))
		return err
	}
	return nil
}

func (i *interactor) DeleteLease(ctx context.Context, nodeId, workloadId string) error {
	logger := i.logger.With(
		zap.String(domain.SDomain, domain.SWorkload),
		zap.String(domain.SNodeID, nodeId),
		zap.String(domain.SWorkloadID, workloadId),
	)

	if err := i.wokrloadRepo.DeleteLease(ctx, nodeId, workloadId); err != nil {
		logger.Error("failed to delete lease", zap.Error(err))
		return err
	}
	logger.Info("deleted lease")
	return nil
}
