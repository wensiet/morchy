package workload

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/wernsiet/morchy/pkg/controlplane/domain"
	"github.com/wernsiet/morchy/pkg/controlplane/domain/workload"
)

func (r *Repository) leasePrimitiveSelect(ctx context.Context, query string, options ...any) (*workload.Lease, error) {
	var l dbLease
	err := r.db.QueryRow(ctx, query, options...).Scan(&l.NodeID, &l.WorkloadID, &l.CreatedAt, &l.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrorWorkloadRepositoryNotFound.New("lease not found")
		}
		return nil, domain.ErrorWorkloadRepositoryInternalError.Wrap(err)
	}
	return l.ToDomain(), nil
}

func (r *Repository) leasePrimitiveExec(ctx context.Context, query string, options ...any) error {
	if _, err := r.db.Exec(ctx, query, options...); err != nil {
		return domain.ErrorWorkloadRepositoryInternalError.Wrap(err)
	}
	return nil
}

func (r *Repository) GetLease(ctx context.Context, nodeId, workloadId string) (*workload.Lease, error) {
	return r.leasePrimitiveSelect(ctx, r.queries.GetLease(), nodeId, workloadId)
}

func (r *Repository) CreateLease(ctx context.Context, nodeId, workloadId string) (*workload.Lease, error) {
	return r.leasePrimitiveSelect(ctx, r.queries.CreateLease(), nodeId, workloadId)
}

func (r *Repository) DeleteExpiredLeases(ctx context.Context, retentionInterval int) error {
	return r.leasePrimitiveExec(ctx, r.queries.DeleteExpiredLeases(), retentionInterval)
}

func (r *Repository) UpdateLease(ctx context.Context, nodeId, workloadId string) error {
	return r.leasePrimitiveExec(ctx, r.queries.UpdateLeaseUpdatedAt(), nodeId, workloadId)
}
