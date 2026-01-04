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

func (r *Repository) ListEvents(ctx context.Context, payloadFilters map[string]string, limit int) ([]*workload.Event, error) {
	query, arguments := r.queries.SelectManyEvents(payloadFilters, limit)
	rows, err := r.db.Query(ctx, query, arguments...)
	if err != nil {
		return nil, domain.ErrorWorkloadRepositoryInternalError.Wrap(err)
	}
	defer rows.Close()

	var events []*workload.Event
	for rows.Next() {
		var e dbEvent

		err := rows.Scan(
			&e.ID,
			&e.SourceID,
			&e.NodeID,
			&e.Payload,
			&e.ProducedAt,
			&e.CreatedAt,
		)
		if err != nil {
			return nil, domain.ErrorWorkloadRepositoryInternalError.Wrap(err)
		}

		events = append(events, e.ToDomain())
	}

	return events, nil
}
