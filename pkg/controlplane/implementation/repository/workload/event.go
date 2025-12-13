package workload

import (
	"context"

	"github.com/wernsiet/morchy/pkg/controlplane/domain"
	"github.com/wernsiet/morchy/pkg/controlplane/domain/workload"
)

func (r *Repository) SaveEvent(ctx context.Context, event workload.Event) error {
	query := r.queries.SaveEvent()
	if _, err := r.db.Exec(
		ctx, query, event.ID, event.SourceID,
		event.NodeID, event.Payload, event.ProducedAt,
	); err != nil {
		return domain.ErrorWorkloadRepositoryInternalError.Wrap(err)
	}
	return nil
}
