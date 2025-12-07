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

func NewWorkloadResponseFromDomain(w *workload.Workload) *WorkloadResponse {
	return &WorkloadResponse{
		ID:        w.ID,
		Status:    string(w.Status),
		Container: w.Spec.Container,
	}
}

type WorkloadSpecRequest struct {
	Container runtime.Container `json:"container"`
}

func (wsr *WorkloadSpecRequest) ToDomain() workload.WorkloadSpec {
	return workload.WorkloadSpec{
		Container: wsr.Container,
	}
}
