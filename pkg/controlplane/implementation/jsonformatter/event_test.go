package jsonformatter

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEventCreateRequest_ToDomain(t *testing.T) {
	t.Run("convert event to domain", func(t *testing.T) {
		payload := json.RawMessage(`{"action":"join","status":"success"}`)
		producedAt := time.Now()

		request := EventCreateRequest{
			ID:         "source-event-id",
			ProducedAt: producedAt,
			Payload:    payload,
		}

		event := request.ToDomain("node-1")

		require.NotEmpty(t, event.ID)
		assert.NotEqual(t, request.ID, event.ID, "Domain event should have new UUID")
		assert.Equal(t, request.ID, event.SourceID)
		assert.Equal(t, "node-1", event.NodeID)
		assert.Equal(t, producedAt, event.ProducedAt)
		assert.Equal(t, payload, event.Payload)
	})

	t.Run("empty payload", func(t *testing.T) {
		request := EventCreateRequest{
			ID:         "event-1",
			ProducedAt: time.Now(),
			Payload:    json.RawMessage{},
		}

		event := request.ToDomain("node-1")

		assert.Equal(t, json.RawMessage{}, event.Payload)
	})

	t.Run("large payload", func(t *testing.T) {
		largeData := make(map[string]interface{})
		for i := 0; i < 1000; i++ {
			largeData[string(rune(i))] = i
		}
		payload, _ := json.Marshal(largeData)

		request := EventCreateRequest{
			ID:         "event-1",
			ProducedAt: time.Now(),
			Payload:    json.RawMessage(payload),
		}

		event := request.ToDomain("node-1")

		assert.Greater(t, len(event.Payload), 1000)
	})

	t.Run("different node IDs", func(t *testing.T) {
		request := EventCreateRequest{
			ID:         "event-1",
			ProducedAt: time.Now(),
			Payload:    json.RawMessage(`{}`),
		}

		event1 := request.ToDomain("node-1")
		event2 := request.ToDomain("node-2")

		assert.Equal(t, "node-1", event1.NodeID)
		assert.Equal(t, "node-2", event2.NodeID)
		assert.NotEqual(t, event1.ID, event2.ID, "Each conversion should create new UUID")
	})

	t.Run("timestamp handling", func(t *testing.T) {
		now := time.Now()

		request := EventCreateRequest{
			ID:         "event-1",
			ProducedAt: now,
			Payload:    json.RawMessage(`{}`),
		}

		event := request.ToDomain("node-1")

		assert.Equal(t, now, event.ProducedAt)
		assert.True(t, event.CreatedAt.IsZero(), "CreatedAt should not be set by ToDomain")
	})
}
