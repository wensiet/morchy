package usecase

import (
	"context"

	"github.com/wernsiet/morchy/internal/domain"
	"github.com/wernsiet/morchy/internal/domain/workload"
)

func (i *interactor) GetLeaseByNodeAndWorkloadID(ctx context.Context, nodeId, workloadId string) (*workload.Lease, error) {
	lease, err := i.wokrloadRepo.GetLease(ctx, nodeId, workloadId)
	if err != nil {
		return nil, err
	}
	return lease, nil
}

func (i *interactor) CreateLease(ctx context.Context, nodeId, workloadId string) (*workload.Lease, error) {
	lease, err := i.wokrloadRepo.CreateLease(ctx, nodeId, workloadId)
	if err != nil {
		return nil, err
	}
	return lease, nil
}

func (i *interactor) ExtendLease(ctx context.Context, nodeId, workloadId string) error {
	err := i.wokrloadRepo.UpdateLease(ctx, nodeId, workloadId)
	if err != nil {
		return nil
	}
	return err
}

func (i *interactor) ExpireLeases(ctx context.Context) error {
	err := i.wokrloadRepo.DeleteExpiredLeases(ctx, domain.CLeaseLifetime)
	if err != nil {
		return nil
	}
	return err
}
