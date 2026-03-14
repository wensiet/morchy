package workload

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWorkloadStatus_Values(t *testing.T) {
	t.Run("all status constants have unique values", func(t *testing.T) {
		statuses := []WorkloadStatus{
			NewWorkloadStatus,
			ActiveWorkloadStatus,
			TerminatedWorkloadStatus,
		}

		uniqueValues := make(map[string]bool)
		for _, status := range statuses {
			value := string(status)
			assert.False(t, uniqueValues[value], "Status value %s is duplicated", value)
			uniqueValues[value] = true
		}

		assert.Equal(t, 3, len(uniqueValues), "Should have 3 unique status values")
	})

	t.Run("status values are as expected", func(t *testing.T) {
		assert.Equal(t, "new", string(NewWorkloadStatus))
		assert.Equal(t, "active", string(ActiveWorkloadStatus))
		assert.Equal(t, "terminated", string(TerminatedWorkloadStatus))
	})
}

func TestWorkload_Fields(t *testing.T) {
	t.Run("workload with all fields", func(t *testing.T) {
		workload := &Workload{
			ID: "workload-1",
			Container: Container{
				Name: "test-container",
			},
		}

		assert.Equal(t, "workload-1", workload.ID)
		assert.Equal(t, "test-container", workload.Container.Name)
	})

	t.Run("workload with minimal fields", func(t *testing.T) {
		workload := &Workload{
			ID: "workload-1",
		}

		assert.Equal(t, "workload-1", workload.ID)
		assert.Empty(t, workload.Container.Name)
	})

	t.Run("workload with empty container name", func(t *testing.T) {
		workload := &Workload{
			ID: "workload-1",
			Container: Container{
				Name: "",
			},
		}

		assert.Empty(t, workload.Container.Name)
	})
}

func TestContainer_Fields(t *testing.T) {
	t.Run("container with name", func(t *testing.T) {
		container := Container{
			Name: "test-container",
		}

		assert.Equal(t, "test-container", container.Name)
	})

	t.Run("container with empty name", func(t *testing.T) {
		container := Container{
			Name: "",
		}

		assert.Empty(t, container.Name)
	})

	t.Run("container name with special characters", func(t *testing.T) {
		specialNames := []string{
			"container-123",
			"container_123",
			"Container123",
			"container.name",
			"container-name-v2",
		}

		for _, name := range specialNames {
			container := Container{Name: name}
			assert.Equal(t, name, container.Name)
		}
	})
}
