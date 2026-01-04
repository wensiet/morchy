package usecase

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/wernsiet/morchy/pkg/controlplane/domain"
	"github.com/wernsiet/morchy/pkg/controlplane/domain/workload"
	"go.uber.org/zap"
)

func newEvent(nodeID string, payload json.RawMessage) workload.Event {
	id := uuid.NewString()
	return workload.Event{
		ID:         id,
		SourceID:   id,
		NodeID:     nodeID,
		Payload:    payload,
		ProducedAt: time.Now(),
	}
}

func (i *interactor) PushEvent(ctx context.Context, event workload.Event) error {
	logger := i.logger.With(
		zap.String(domain.SDomain, domain.SWorkload),
		zap.String(domain.SEventSourceID, event.SourceID),
		zap.String(domain.SNodeID, event.NodeID),
	)
	err := i.wokrloadRepo.SaveEvent(ctx, event)
	if err != nil {
		logger.Error("failed to save event", zap.Error(err))
		return err
	}
	return nil
}
