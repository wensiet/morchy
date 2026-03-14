package jsonformatter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wernsiet/morchy/pkg/controlplane/domain/workload"
)

func TestNewLeaseResponseFromDomain(t *testing.T) {
	t.Run("lease with all fields", func(t *testing.T) {
		now := time.Now()
		domainLease := &workload.Lease{
			NodeID:     "node-1",
			WorkloadID: "workload-1",
			CreatedAt:  now,
			UpdatedAt:  now.Add(5 * time.Minute),
		}

		response := NewLeaseResponseFromDomain(domainLease)

		require.NotNil(t, response)
		assert.Equal(t, "node-1", response.NodeID)
		assert.Equal(t, "workload-1", response.WorkloadID)
		assert.Equal(t, now, response.CreatedAt)
		assert.Equal(t, now.Add(5*time.Minute), response.UpdatedAt)
	})

	t.Run("lease with zero timestamps", func(t *testing.T) {
		domainLease := &workload.Lease{
			NodeID:     "node-1",
			WorkloadID: "workload-1",
			CreatedAt:  time.Time{},
			UpdatedAt:  time.Time{},
		}

		response := NewLeaseResponseFromDomain(domainLease)

		require.NotNil(t, response)
		assert.True(t, response.CreatedAt.IsZero())
		assert.True(t, response.UpdatedAt.IsZero())
	})

	t.Run("lease with same created and updated timestamps", func(t *testing.T) {
		now := time.Now()
		domainLease := &workload.Lease{
			NodeID:     "node-1",
			WorkloadID: "workload-1",
			CreatedAt:  now,
			UpdatedAt:  now,
		}

		response := NewLeaseResponseFromDomain(domainLease)

		require.NotNil(t, response)
		assert.Equal(t, now, response.CreatedAt)
		assert.Equal(t, now, response.UpdatedAt)
	})

	t.Run("lease with updated before created (unlikely but possible)", func(t *testing.T) {
		now := time.Now()
		domainLease := &workload.Lease{
			NodeID:     "node-1",
			WorkloadID: "workload-1",
			CreatedAt:  now,
			UpdatedAt:  now.Add(-1 * time.Second),
		}

		response := NewLeaseResponseFromDomain(domainLease)

		require.NotNil(t, response)
		assert.True(t, response.UpdatedAt.Before(response.CreatedAt))
	})
}
