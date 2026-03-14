package workload

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWorkloadStatus_Values(t *testing.T) {
	t.Run("all status constants have unique values", func(t *testing.T) {
		statuses := []WorkloadStatus{
			NewWorkloadStatus,
			PendingWorkloadStatus,
			StuckWorkloadStatus,
			ActiveWorkloadStatus,
			FailedWorkloadStatus,
			DegradedWorkloadStatus,
		}

		uniqueValues := make(map[string]bool)
		for _, status := range statuses {
			value := string(status)
			assert.False(t, uniqueValues[value], "Status value %s is duplicated", value)
			uniqueValues[value] = true
		}

		assert.Equal(t, 6, len(uniqueValues), "Should have 6 unique status values")
	})

	t.Run("status values are as expected", func(t *testing.T) {
		assert.Equal(t, "new", string(NewWorkloadStatus))
		assert.Equal(t, "pending", string(PendingWorkloadStatus))
		assert.Equal(t, "stuck", string(StuckWorkloadStatus))
		assert.Equal(t, "active", string(ActiveWorkloadStatus))
		assert.Equal(t, "failed", string(FailedWorkloadStatus))
		assert.Equal(t, "degraded", string(DegradedWorkloadStatus))
	})
}

func TestWorkloadSpec_Ports(t *testing.T) {
	t.Run("spec with ports", func(t *testing.T) {
		containerPort := 8080
		hostPort := 80

		spec := WorkloadSpec{
			Name:          "test",
			Image:         "nginx",
			CPU:           100,
			RAM:           256,
			ContainerPort: &containerPort,
			HostPort:      &hostPort,
		}

		assert.NotNil(t, spec.ContainerPort)
		assert.NotNil(t, spec.HostPort)
		assert.Equal(t, 8080, *spec.ContainerPort)
		assert.Equal(t, 80, *spec.HostPort)
	})

	t.Run("spec without ports", func(t *testing.T) {
		spec := WorkloadSpec{
			Name:  "test",
			Image: "nginx",
			CPU:   100,
			RAM:   256,
		}

		assert.Nil(t, spec.ContainerPort)
		assert.Nil(t, spec.HostPort)
	})
}

func TestWorkloadSpec_Command(t *testing.T) {
	t.Run("spec with command", func(t *testing.T) {
		spec := WorkloadSpec{
			Name:    "test",
			Image:   "nginx",
			CPU:     100,
			RAM:     256,
			Command: []string{"nginx", "-g", "daemon off;"},
		}

		assert.NotNil(t, spec.Command)
		assert.Equal(t, 3, len(spec.Command))
		assert.Equal(t, "nginx", spec.Command[0])
	})

	t.Run("spec without command", func(t *testing.T) {
		spec := WorkloadSpec{
			Name:  "test",
			Image: "nginx",
			CPU:   100,
			RAM:   256,
		}

		assert.Nil(t, spec.Command)
	})
}

func TestWorkloadSpec_Env(t *testing.T) {
	t.Run("spec with environment variables", func(t *testing.T) {
		spec := WorkloadSpec{
			Name:  "test",
			Image: "nginx",
			CPU:   100,
			RAM:   256,
			Env: map[string]string{
				"ENV1": "value1",
				"ENV2": "value2",
			},
		}

		assert.NotNil(t, spec.Env)
		assert.Equal(t, 2, len(spec.Env))
		assert.Equal(t, "value1", spec.Env["ENV1"])
		assert.Equal(t, "value2", spec.Env["ENV2"])
	})

	t.Run("spec without environment variables", func(t *testing.T) {
		spec := WorkloadSpec{
			Name:  "test",
			Image: "nginx",
			CPU:   100,
			RAM:   256,
		}

		assert.Nil(t, spec.Env)
	})
}

func TestWorkload_Fields(t *testing.T) {
	t.Run("workload with all fields", func(t *testing.T) {
		containerPort := 8080
		hostPort := 80

		w := &Workload{
			ID:     "test-id",
			Status: ActiveWorkloadStatus,
			Spec: WorkloadSpec{
				Name:          "test-workload",
				Image:         "nginx:latest",
				CPU:           100,
				RAM:           256,
				Command:       []string{"nginx"},
				Env:           map[string]string{"ENV1": "value1"},
				ContainerPort: &containerPort,
				HostPort:      &hostPort,
			},
		}

		assert.Equal(t, "test-id", w.ID)
		assert.Equal(t, ActiveWorkloadStatus, w.Status)
		assert.Equal(t, "test-workload", w.Spec.Name)
		assert.Equal(t, "nginx:latest", w.Spec.Image)
		assert.Equal(t, uint(100), w.Spec.CPU)
		assert.Equal(t, uint(256), w.Spec.RAM)
		assert.Equal(t, []string{"nginx"}, w.Spec.Command)
		assert.Equal(t, "value1", w.Spec.Env["ENV1"])
		assert.Equal(t, 8080, *w.Spec.ContainerPort)
		assert.Equal(t, 80, *w.Spec.HostPort)
	})

	t.Run("workload with minimal fields", func(t *testing.T) {
		w := &Workload{
			ID:     "test-id",
			Status: NewWorkloadStatus,
			Spec: WorkloadSpec{
				Name:  "test-workload",
				Image: "nginx:latest",
				CPU:   100,
				RAM:   256,
			},
		}

		assert.Equal(t, "test-id", w.ID)
		assert.Equal(t, NewWorkloadStatus, w.Status)
		assert.Equal(t, "test-workload", w.Spec.Name)
		assert.Equal(t, "nginx:latest", w.Spec.Image)
		assert.Nil(t, w.Spec.Command)
		assert.Nil(t, w.Spec.Env)
		assert.Nil(t, w.Spec.ContainerPort)
		assert.Nil(t, w.Spec.HostPort)
	})
}
