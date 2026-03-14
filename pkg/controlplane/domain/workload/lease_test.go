package workload

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLease_Fields(t *testing.T) {
	t.Run("lease with all fields", func(t *testing.T) {
		now := time.Now()
		lease := &Lease{
			NodeID:     "node-1",
			WorkloadID: "workload-1",
			CreatedAt:  now,
			UpdatedAt:  now.Add(5 * time.Minute),
		}

		assert.Equal(t, "node-1", lease.NodeID)
		assert.Equal(t, "workload-1", lease.WorkloadID)
		assert.Equal(t, now, lease.CreatedAt)
		assert.Equal(t, now.Add(5*time.Minute), lease.UpdatedAt)
	})

	t.Run("lease with zero timestamps", func(t *testing.T) {
		lease := &Lease{
			NodeID:     "node-1",
			WorkloadID: "workload-1",
			CreatedAt:  time.Time{},
			UpdatedAt:  time.Time{},
		}

		assert.Equal(t, "node-1", lease.NodeID)
		assert.Equal(t, "workload-1", lease.WorkloadID)
		assert.True(t, lease.CreatedAt.IsZero())
		assert.True(t, lease.UpdatedAt.IsZero())
	})
}

func TestLease_TimeComparison(t *testing.T) {
	t.Run("compare created and updated timestamps", func(t *testing.T) {
		now := time.Now()
		lease := &Lease{
			NodeID:     "node-1",
			WorkloadID: "workload-1",
			CreatedAt:  now,
			UpdatedAt:  now.Add(30 * time.Second),
		}

		assert.True(t, lease.UpdatedAt.After(lease.CreatedAt))
		assert.False(t, lease.CreatedAt.After(lease.UpdatedAt))
	})

	t.Run("check if lease is stale", func(t *testing.T) {
		staleLease := &Lease{
			NodeID:     "node-1",
			WorkloadID: "workload-1",
			CreatedAt:  time.Now().Add(-1 * time.Hour),
			UpdatedAt:  time.Now().Add(-1 * time.Hour),
		}

		timeSinceUpdate := time.Since(staleLease.UpdatedAt)
		assert.True(t, timeSinceUpdate > 30*time.Second)
	})

	t.Run("check if lease is fresh", func(t *testing.T) {
		freshLease := &Lease{
			NodeID:     "node-1",
			WorkloadID: "workload-1",
			CreatedAt:  time.Now().Add(-1 * time.Hour),
			UpdatedAt:  time.Now().Add(-10 * time.Second),
		}

		timeSinceUpdate := time.Since(freshLease.UpdatedAt)
		assert.True(t, timeSinceUpdate < 30*time.Second)
	})
}

func TestLease_IDs(t *testing.T) {
	t.Run("empty node ID", func(t *testing.T) {
		lease := &Lease{
			NodeID:     "",
			WorkloadID: "workload-1",
		}

		assert.Equal(t, "", lease.NodeID)
		assert.False(t, lease.NodeID != "")
	})

	t.Run("empty workload ID", func(t *testing.T) {
		lease := &Lease{
			NodeID:     "node-1",
			WorkloadID: "",
		}

		assert.Equal(t, "", lease.WorkloadID)
		assert.False(t, lease.WorkloadID != "")
	})

	t.Run("both IDs populated", func(t *testing.T) {
		lease := &Lease{
			NodeID:     "node-1",
			WorkloadID: "workload-1",
		}

		assert.NotEmpty(t, lease.NodeID)
		assert.NotEmpty(t, lease.WorkloadID)
	})
}
