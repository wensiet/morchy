package workload

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/wernsiet/morchy/pkg/controlplane/domain"
	"github.com/wernsiet/morchy/pkg/controlplane/domain/workload"
	"github.com/wernsiet/morchy/pkg/runtime"
)

func (r *Repository) workloadPrimitiveSelect(ctx context.Context, query string, options ...any) (*workload.Workload, error) {
	var w dbWorkload
	var s dbWorkloadSpec

	err := r.db.QueryRow(ctx, query, options...).Scan(
		&w.ID,
		&w.Status,
		&w.CreatedAt,
		&s.ID,
		&s.Image,
		&s.CPU,
		&s.RAM,
		&s.Command,
		&s.Env,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrorWorkloadRepositoryNotFound.New("workload not found")
		}
		return nil, domain.ErrorWorkloadRepositoryInternalError.Wrap(err)
	}

	w.Spec = s
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
		var s dbWorkloadSpec

		err := rows.Scan(
			&w.ID,
			&w.Status,
			&w.CreatedAt,
			&s.ID,
			&s.Image,
			&s.CPU,
			&s.RAM,
			&s.Command,
			&s.Env,
			nil,
		)
		if err != nil {
			return nil, domain.ErrorWorkloadRepositoryInternalError.Wrap(err)
		}

		w.Spec = s
		workloads = append(workloads, w.ToDomain())
	}

	return workloads, nil
}

func (r *Repository) CreateWorkload(ctx context.Context, domainWorkload workload.Workload) (*workload.Workload, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, domain.ErrorWorkloadRepositoryInternalError.Wrap(err)
	}
	defer tx.Rollback(ctx)

	var w dbWorkload
	err = tx.QueryRow(
		ctx,
		r.queries.CreateWorkload(),
		domainWorkload.ID,
		domainWorkload.Status,
	).Scan(&w.ID, &w.Status, &w.CreatedAt)
	if err != nil {
		return nil, domain.ErrorWorkloadRepositoryInternalError.Wrap(err)
	}

	var s dbWorkloadSpec
	err = tx.QueryRow(
		ctx,
		r.queries.CreateWorkloadSpec(),
		w.ID,
		domainWorkload.Spec.Image,
		domainWorkload.Spec.CPU,
		domainWorkload.Spec.RAM,
		domainWorkload.Spec.Command,
		domainWorkload.Spec.Env,
	).Scan(
		&s.ID,
		&s.Image,
		&s.CPU,
		&s.RAM,
		&s.Command,
		&s.Env,
	)
	if err != nil {
		return nil, domain.ErrorWorkloadRepositoryInternalError.Wrap(err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, domain.ErrorWorkloadRepositoryInternalError.Wrap(err)
	}

	w.Spec = s
	return w.ToDomain(), nil
}

func (r *Repository) GetWorkload(ctx context.Context, workloadID string) (*workload.Workload, error) {
	return r.workloadPrimitiveSelect(ctx, r.queries.SelectWorkloadByID(), workloadID)
}
