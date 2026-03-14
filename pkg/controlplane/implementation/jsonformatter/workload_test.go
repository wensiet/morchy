package jsonformatter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wernsiet/morchy/pkg/controlplane/domain/workload"
)

func TestMapEnvs(t *testing.T) {
	t.Run("map environment variables", func(t *testing.T) {
		envMap := map[string]string{
			"ENV1": "value1",
			"ENV2": "value2",
			"ENV3": "value3",
		}

		result := mapEnvs(envMap)

		require.Equal(t, 3, len(result))

		resultMap := make(map[string]string)
		for _, env := range result {
			resultMap[env.Name] = env.Value
		}

		assert.Equal(t, "value1", resultMap["ENV1"])
		assert.Equal(t, "value2", resultMap["ENV2"])
		assert.Equal(t, "value3", resultMap["ENV3"])
	})

	t.Run("empty environment map", func(t *testing.T) {
		envMap := map[string]string{}

		result := mapEnvs(envMap)

		assert.Equal(t, 0, len(result))
	})

	t.Run("nil environment map", func(t *testing.T) {
		var envMap map[string]string = nil

		result := mapEnvs(envMap)

		assert.Equal(t, 0, len(result))
	})
}

func TestNewWorkloadResponseFromDomain(t *testing.T) {
	t.Run("workload with all fields", func(t *testing.T) {
		containerPort := 8080
		hostPort := 80

		domainWorkload := &workload.Workload{
			ID:     "workload-1",
			Status: workload.ActiveWorkloadStatus,
			Spec: workload.WorkloadSpec{
				Name:          "test-workload",
				Image:         "nginx:latest",
				CPU:           100,
				RAM:           256,
				Command:       []string{"nginx", "-g", "daemon off;"},
				Env:           map[string]string{"ENV1": "value1"},
				ContainerPort: &containerPort,
				HostPort:      &hostPort,
			},
		}

		response := NewWorkloadResponseFromDomain(domainWorkload)

		require.NotNil(t, response)
		assert.Equal(t, "workload-1", response.ID)
		assert.Equal(t, "active", response.Status)
		assert.Equal(t, "test-workload", response.Container.Name)
		assert.Equal(t, "nginx:latest", response.Container.Image)
		assert.Equal(t, uint(100), response.Container.Resources.CPU)
		assert.Equal(t, uint(256), response.Container.Resources.RAM)
		assert.Equal(t, []string{"nginx", "-g", "daemon off;"}, response.Container.Command)
		assert.Len(t, response.Container.Env, 1)
		assert.Equal(t, "ENV1", response.Container.Env[0].Name)
		assert.Equal(t, "value1", response.Container.Env[0].Value)
		require.NotNil(t, response.Container.NetConfig)
		assert.Equal(t, 8080, response.Container.NetConfig.ContainerPort)
		assert.Equal(t, 80, response.Container.NetConfig.HostPort)
		assert.Equal(t, "tcp", response.Container.NetConfig.Protocol)
	})

	t.Run("workload without network config", func(t *testing.T) {
		domainWorkload := &workload.Workload{
			ID:     "workload-1",
			Status: workload.NewWorkloadStatus,
			Spec: workload.WorkloadSpec{
				Name:    "test-workload",
				Image:   "nginx:latest",
				CPU:     100,
				RAM:     256,
				Command: []string{"nginx"},
				Env:     map[string]string{},
			},
		}

		response := NewWorkloadResponseFromDomain(domainWorkload)

		require.NotNil(t, response)
		assert.Equal(t, "workload-1", response.ID)
		assert.Equal(t, "new", response.Status)
		assert.Nil(t, response.Container.NetConfig)
	})

	t.Run("workload with only container port", func(t *testing.T) {
		containerPort := 8080

		domainWorkload := &workload.Workload{
			ID:     "workload-1",
			Status: workload.PendingWorkloadStatus,
			Spec: workload.WorkloadSpec{
				Name:          "test-workload",
				Image:         "nginx:latest",
				CPU:           100,
				RAM:           256,
				ContainerPort: &containerPort,
			},
		}

		response := NewWorkloadResponseFromDomain(domainWorkload)

		require.NotNil(t, response)
		assert.Nil(t, response.Container.NetConfig)
	})
}

func TestWorkloadSpecRequest_ToDomain(t *testing.T) {
	t.Run("spec with all fields", func(t *testing.T) {
		containerPort := 8080

		request := &WorkloadSpecRequest{
			Image:         "nginx:latest",
			CPU:           100,
			RAM:           256,
			Command:       []string{"nginx", "-g", "daemon off;"},
			Env:           map[string]string{"ENV1": "value1"},
			ContainerPort: &containerPort,
		}

		spec := request.ToDomain()

		assert.Equal(t, "nginx:latest", spec.Image)
		assert.Equal(t, uint(100), spec.CPU)
		assert.Equal(t, uint(256), spec.RAM)
		assert.Equal(t, []string{"nginx", "-g", "daemon off;"}, spec.Command)
		assert.Equal(t, "value1", spec.Env["ENV1"])
		assert.NotNil(t, spec.ContainerPort)
		assert.Equal(t, 8080, *spec.ContainerPort)
		assert.NotNil(t, spec.HostPort)
		assert.GreaterOrEqual(t, *spec.HostPort, 32766)
		assert.LessOrEqual(t, *spec.HostPort, 65534)
	})

	t.Run("spec without container port", func(t *testing.T) {
		request := &WorkloadSpecRequest{
			Image:   "nginx:latest",
			CPU:     100,
			RAM:     256,
			Command: []string{"nginx"},
			Env:     map[string]string{},
		}

		spec := request.ToDomain()

		assert.Equal(t, "nginx:latest", spec.Image)
		assert.Nil(t, spec.ContainerPort)
		assert.Nil(t, spec.HostPort)
	})

	t.Run("spec with empty fields", func(t *testing.T) {
		request := &WorkloadSpecRequest{
			Image: "nginx:latest",
			CPU:   0,
			RAM:   0,
		}

		spec := request.ToDomain()

		assert.Equal(t, "nginx:latest", spec.Image)
		assert.Equal(t, uint(0), spec.CPU)
		assert.Equal(t, uint(0), spec.RAM)
		assert.Nil(t, spec.Command)
		assert.Nil(t, spec.Env)
	})
}

func TestGetRandomPort(t *testing.T) {
	t.Run("generates port in valid range", func(t *testing.T) {
		port := getRandomPort()

		assert.GreaterOrEqual(t, port, 32766)
		assert.LessOrEqual(t, port, 65534)
	})

	t.Run("generates different ports", func(t *testing.T) {
		ports := make(map[int]bool)

		for i := 0; i < 100; i++ {
			port := getRandomPort()
			ports[port] = true
		}

		assert.Greater(t, len(ports), 1, "Should generate different ports")
	})
}
