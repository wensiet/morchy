package workload

import (
	"context"

	"github.com/wernsiet/morchy/pkg/controlplane/domain"
	"github.com/wernsiet/morchy/pkg/controlplane/domain/workload"
)

func (r *Repository) ListEdges(ctx context.Context) ([]*workload.Edge, error) {
	rows, err := r.db.Query(ctx, r.queries.SelectManyEdges())
	if err != nil {
		return nil, domain.ErrorWorkloadRepositoryInternalError.Wrap(err)
	}
	defer rows.Close()

	var edges []*workload.Edge
	for rows.Next() {
		var e dbEdge
		err = rows.Scan(
			&e.WorkloadID,
			&e.NodeID,
			&e.HostPort,
		)
		if err != nil {
			return nil, domain.ErrorWorkloadRepositoryInternalError.Wrap(err)
		}
		edges = append(edges, e.ToDomain())
	}

	return edges, nil
}
