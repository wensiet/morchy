package workload

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wernsiet/morchy/pkg/agent/domain"
)

func TestEventPayload_Fields(t *testing.T) {
	t.Run("event payload with all fields", func(t *testing.T) {
		payload := EventPayload{
			WorkloadID: "workload-1",
			Action:     "healthcheck",
			Status:     domain.EventActionStatusSuccess,
		}

		assert.Equal(t, "workload-1", payload.WorkloadID)
		assert.Equal(t, "healthcheck", payload.Action)
		assert.Equal(t, domain.EventActionStatusSuccess, payload.Status)
	})

	t.Run("event payload with different statuses", func(t *testing.T) {
		tests := []struct {
			name   string
			status domain.EventActionStatus
		}{
			{"success", domain.EventActionStatusSuccess},
			{"failure", domain.EventActionStatusFailed},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				payload := EventPayload{
					WorkloadID: "workload-1",
					Action:     "action",
					Status:     tt.status,
				}

				assert.Equal(t, tt.status, payload.Status)
			})
		}
	})

	t.Run("event payload with various actions", func(t *testing.T) {
		actions := []string{
			"healthcheck",
			"start",
			"stop",
			"create",
			"delete",
		}

		for _, action := range actions {
			payload := EventPayload{
				WorkloadID: "workload-1",
				Action:     action,
				Status:     domain.EventActionStatusSuccess,
			}

			assert.Equal(t, action, payload.Action)
		}
	})

	t.Run("event payload with empty fields", func(t *testing.T) {
		payload := EventPayload{}

		assert.Empty(t, payload.WorkloadID)
		assert.Empty(t, payload.Action)
		assert.Empty(t, string(payload.Status))
	})
}
