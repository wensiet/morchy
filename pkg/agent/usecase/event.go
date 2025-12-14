package usecase

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/wernsiet/morchy/pkg/agent/domain"
	"github.com/wernsiet/morchy/pkg/agent/domain/workload"
	"go.uber.org/zap"

	apitypes "github.com/wernsiet/morchy/pkg/controlplane/implementation/jsonformatter"
)

func (i *interactor) PushEvent(ctx context.Context, payload workload.EventPayload) error {
	payloadJson, err := json.Marshal(payload)
	if err != nil {
		return domain.ErrorBaseWorkloadInternal.Wrapf(err, "failed to marshal event payload to JSON")
	}
	return i.controlplaneClient.PushEvent(ctx, apitypes.EventCreateRequest{
		ID:         uuid.NewString(),
		ProducedAt: time.Now(),
		Payload:    payloadJson,
	})
}

func (i *interactor) asyncPushEvent(ctx context.Context, payload workload.EventPayload) {
	go func() {
		if err := i.PushEvent(ctx, payload); err != nil {
			i.logger.Warn("failed to push event", zap.Error(err))
		}
	}()
}
