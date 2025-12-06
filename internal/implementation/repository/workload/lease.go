package workload

import (
	"context"

	"github.com/wernsiet/morchy/internal/domain/workload"
)

func (r *Repository) leasePrimitiveSelect(ctx context.Context, query string, options ...any) (*workload.Lease, error) {
	var l dbLease
	err := r.db.QueryRow(ctx, query, options...).Scan(&l.NodeID, &l.WorkloadID, &l.CreatedAt, &l.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return l.ToDomain(), nil
}

func (r *Repository) leasePrimitiveExec(ctx context.Context, query string, options ...any) error {
	if _, err := r.db.Exec(ctx, query, options...); err != nil {
		return err
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
