package workload

import (
	"encoding/json"
	"time"
)

type Event struct {
	ID         string
	SourceID   string
	NodeID     string
	Payload    json.RawMessage
	ProducedAt time.Time
	CreatedAt  time.Time
}
