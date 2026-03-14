package workload

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEvent_Fields(t *testing.T) {
	t.Run("event with all fields", func(t *testing.T) {
		now := time.Now()
		payload := json.RawMessage(`{"status":"success","message":"test"}`)

		event := &Event{
			ID:         "event-1",
			SourceID:   "source-1",
			NodeID:     "node-1",
			Payload:    payload,
			ProducedAt: now,
			CreatedAt:  now.Add(1 * time.Second),
		}

		assert.Equal(t, "event-1", event.ID)
		assert.Equal(t, "source-1", event.SourceID)
		assert.Equal(t, "node-1", event.NodeID)
		assert.Equal(t, payload, event.Payload)
		assert.Equal(t, now, event.ProducedAt)
		assert.Equal(t, now.Add(1*time.Second), event.CreatedAt)
	})

	t.Run("event with minimal fields", func(t *testing.T) {
		event := &Event{
			ID:       "event-1",
			SourceID: "source-1",
			NodeID:   "node-1",
			Payload:  json.RawMessage(`{}`),
		}

		assert.Equal(t, "event-1", event.ID)
		assert.Equal(t, "source-1", event.SourceID)
		assert.Equal(t, "node-1", event.NodeID)
		assert.Equal(t, json.RawMessage(`{}`), event.Payload)
		assert.True(t, event.ProducedAt.IsZero())
		assert.True(t, event.CreatedAt.IsZero())
	})
}

func TestEvent_Payload(t *testing.T) {
	t.Run("event with JSON payload", func(t *testing.T) {
		payloadData := map[string]interface{}{
			"status":  "success",
			"message": "test message",
			"count":   42,
		}
		payload, _ := json.Marshal(payloadData)

		event := &Event{
			ID:       "event-1",
			SourceID: "source-1",
			NodeID:   "node-1",
			Payload:  json.RawMessage(payload),
		}

		assert.NotNil(t, event.Payload)
		assert.Equal(t, json.RawMessage(payload), event.Payload)

		var decoded map[string]interface{}
		err := json.Unmarshal(event.Payload, &decoded)
		require.NoError(t, err)
		assert.Equal(t, "success", decoded["status"])
		assert.Equal(t, "test message", decoded["message"])
		assert.Equal(t, float64(42), decoded["count"])
	})

	t.Run("event with empty payload", func(t *testing.T) {
		event := &Event{
			ID:       "event-1",
			SourceID: "source-1",
			NodeID:   "node-1",
			Payload:  json.RawMessage{},
		}

		assert.Equal(t, json.RawMessage{}, event.Payload)
	})

	t.Run("event with large payload", func(t *testing.T) {
		largeData := make(map[string]interface{})
		for i := 0; i < 100; i++ {
			largeData[string(rune(i))] = i
		}
		payload, _ := json.Marshal(largeData)

		event := &Event{
			ID:       "event-1",
			SourceID: "source-1",
			NodeID:   "node-1",
			Payload:  json.RawMessage(payload),
		}

		assert.Greater(t, len(event.Payload), 100)
	})
}

func TestEvent_Timestamps(t *testing.T) {
	t.Run("produced before created", func(t *testing.T) {
		now := time.Now()
		event := &Event{
			ID:         "event-1",
			SourceID:   "source-1",
			NodeID:     "node-1",
			Payload:    json.RawMessage(`{}`),
			ProducedAt: now,
			CreatedAt:  now.Add(1 * time.Second),
		}

		assert.True(t, event.CreatedAt.After(event.ProducedAt))
	})

	t.Run("produced and created at same time", func(t *testing.T) {
		now := time.Now()
		event := &Event{
			ID:         "event-1",
			SourceID:   "source-1",
			NodeID:     "node-1",
			Payload:    json.RawMessage(`{}`),
			ProducedAt: now,
			CreatedAt:  now,
		}

		assert.Equal(t, event.ProducedAt, event.CreatedAt)
	})

	t.Run("zero timestamps", func(t *testing.T) {
		event := &Event{
			ID:         "event-1",
			SourceID:   "source-1",
			NodeID:     "node-1",
			Payload:    json.RawMessage(`{}`),
			ProducedAt: time.Time{},
			CreatedAt:  time.Time{},
		}

		assert.True(t, event.ProducedAt.IsZero())
		assert.True(t, event.CreatedAt.IsZero())
	})
}

func TestEvent_IDs(t *testing.T) {
	t.Run("all IDs populated", func(t *testing.T) {
		event := &Event{
			ID:       "event-1",
			SourceID: "source-1",
			NodeID:   "node-1",
			Payload:  json.RawMessage(`{}`),
		}

		assert.NotEmpty(t, event.ID)
		assert.NotEmpty(t, event.SourceID)
		assert.NotEmpty(t, event.NodeID)
	})

	t.Run("empty IDs", func(t *testing.T) {
		event := &Event{
			ID:       "",
			SourceID: "",
			NodeID:   "",
			Payload:  json.RawMessage(`{}`),
		}

		assert.Empty(t, event.ID)
		assert.Empty(t, event.SourceID)
		assert.Empty(t, event.NodeID)
	})

	t.Run("partial IDs", func(t *testing.T) {
		event := &Event{
			ID:       "event-1",
			SourceID: "",
			NodeID:   "node-1",
			Payload:  json.RawMessage(`{}`),
		}

		assert.NotEmpty(t, event.ID)
		assert.Empty(t, event.SourceID)
		assert.NotEmpty(t, event.NodeID)
	})
}
