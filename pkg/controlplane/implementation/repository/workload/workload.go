package workload

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/wernsiet/morchy/pkg/controlplane/domain"
	"github.com/wernsiet/morchy/pkg/controlplane/domain/workload"
	"github.com/wernsiet/morchy/pkg/runtime"
)

func (r *Repository) workloadPrimitiveSelect(ctx context.Context, query string, options ...any) (*workload.Workload, error) {
	var w dbWorkload
	err := r.db.QueryRow(ctx, query, options...).Scan(&w.ID, &w.Status, &w.CreatedAt, &w.Container)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrorWorkloadRepositoryNotFound.New("workload not found")
		}
		return nil, domain.ErrorWorkloadRepositoryInternalError.Wrap(err)
	}
	return w.ToDomain(), nil
}

func (r *Repository) ListWorkloads(ctx context.Context, status *string, resources *runtime.ResourceLimits) ([]*workload.Workload, error) {
	query, arguments := r.queries.SelectManyWorkloads(status, resources)
	rows, err := r.db.Query(ctx, query, arguments...)
	if err != nil {
		return nil, domain.ErrorWorkloadRepositoryInternalError.Wrap(err)
	}
	defer rows.Close()

	var workloads []*workload.Workload
	for rows.Next() {
		var w dbWorkload
		err := rows.Scan(&w.ID, &w.Status, &w.CreatedAt, &w.Container, nil)
		if err != nil {
			return nil, domain.ErrorWorkloadRepositoryInternalError.Wrap(err)
		}
		workloads = append(workloads, w.ToDomain())
	}

	return workloads, nil
}

func (r *Repository) CreateWorklod(ctx context.Context, domainWorkload workload.Workload) (*workload.Workload, error) {
	containerJSON, err := json.Marshal(domainWorkload.Spec.Container)
	if err != nil {
		return nil, domain.ErrorWorkloadRepositoryInternalError.Wrap(err)
	}
	return r.workloadPrimitiveSelect(ctx, r.queries.CreateWorkload(), domainWorkload.ID, domainWorkload.Status, containerJSON)
}

func (r *Repository) GetWorkload(ctx context.Context, workloadID string) (*workload.Workload, error) {
	return r.workloadPrimitiveSelect(ctx, r.queries.SelectWorkloadByID(), workloadID)
}
