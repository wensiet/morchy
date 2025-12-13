package jsonformatter

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/wernsiet/morchy/pkg/controlplane/domain/workload"
)

type EventCreateRequest struct {
	ID         string          `json:"id" example:"some-event-id"`
	ProducedAt time.Time       `json:"produced_at" example:"2025-11-03T13:45:00+03"`
	Payload    json.RawMessage `json:"payload" example:"{\"action\": \"join\"}"`
}

func (e EventCreateRequest) ToDomain(nodeID string) workload.Event {
	return workload.Event{
		ID:         uuid.NewString(),
		SourceID:   e.ID,
		NodeID:     nodeID,
		ProducedAt: e.ProducedAt,
		Payload:    e.Payload,
	}
}
