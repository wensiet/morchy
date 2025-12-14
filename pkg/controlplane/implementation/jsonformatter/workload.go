package jsonformatter

import (
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
	return &WorkloadResponse{
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
}

type WorkloadSpecRequest struct {
	Image   string            `json:"image"`
	CPU     uint              `json:"cpu"`
	RAM     uint              `json:"ram"`
	Command []string          `json:"command"`
	Env     map[string]string `json:"env"`
}

func (wsr *WorkloadSpecRequest) ToDomain() workload.WorkloadSpec {
	return workload.WorkloadSpec{
		Image:   wsr.Image,
		Command: wsr.Command,
		Env:     wsr.Env,
		CPU:     wsr.CPU,
		RAM:     wsr.RAM,
	}
}
