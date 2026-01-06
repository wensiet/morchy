package jsonformatter

import (
	"math/rand"
	"time"

	"github.com/wernsiet/morchy/pkg/controlplane/domain/workload"
	"github.com/wernsiet/morchy/pkg/runtime"
)

type WorkloadResponse struct {
	ID        string            `json:"id" example:"some-uuid"`
	Status    string            `json:"status" example:"new"`
	Container runtime.Container `json:"container"`
}

func mapEnvs(env map[string]string) []runtime.EnvVar {
	var list []runtime.EnvVar
	for k, v := range env {
		list = append(list, runtime.EnvVar{
			Name:  k,
			Value: v,
		})
	}
	return list
}

func NewWorkloadResponseFromDomain(w *workload.Workload) *WorkloadResponse {
	wl := WorkloadResponse{
		ID:     w.ID,
		Status: string(w.Status),
		Container: runtime.Container{
			Name:    w.Spec.Name,
			Image:   w.Spec.Image,
			Command: w.Spec.Command,
			Env:     mapEnvs(w.Spec.Env),
			Resources: runtime.ResourceLimits{
				CPU: w.Spec.CPU,
				RAM: w.Spec.RAM,
			},
		},
	}
	if w.Spec.ContainerPort != nil && w.Spec.HostPort != nil {
		wl.Container.NetConfig = &runtime.NetConfig{
			ContainerPort: *w.Spec.ContainerPort,
			HostPort:      *w.Spec.HostPort,
			Protocol:      "tcp",
		}
	}
	return &wl
}

type WorkloadSpecRequest struct {
	Image         string            `json:"image"`
	CPU           uint              `json:"cpu"`
	RAM           uint              `json:"ram"`
	Command       []string          `json:"command"`
	Env           map[string]string `json:"env"`
	ContainerPort *int              `json:"container_port"`
}

// getRandomPort - generates random port from 32766 to 65534
func getRandomPort() int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	return r.Intn(65534-32766+1) + 32766
}

func (wsr *WorkloadSpecRequest) ToDomain() workload.WorkloadSpec {
	wls := workload.WorkloadSpec{
		Image:   wsr.Image,
		Command: wsr.Command,
		Env:     wsr.Env,
		CPU:     wsr.CPU,
		RAM:     wsr.RAM,
	}
	if wsr.ContainerPort != nil {
		wls.ContainerPort = wsr.ContainerPort
		hostPort := getRandomPort()
		wls.HostPort = &hostPort
	}
	return wls
}
