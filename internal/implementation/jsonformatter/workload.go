package jsonformatter

import (
	"github.com/wernsiet/morchy/internal/domain/workload"
	pkgworkload "github.com/wernsiet/morchy/pkg/workload"
)

type WorkloadResponse struct {
	ID     string `json:"id" example:"some-uuid"`
	Status string `json:"status" example:"new"`
}

func NewWorkloadResponseFromDomain(w *workload.Workload) *WorkloadResponse {
	return &WorkloadResponse{
		ID:     w.ID,
		Status: string(w.Status),
	}
}

type WorkloadSpecRequest struct {
	Container pkgworkload.Container `json:"container"`
}

func (wsr *WorkloadSpecRequest) ToDomain() workload.WorkloadSpec {
	return workload.WorkloadSpec{
		Container: wsr.Container,
	}
}
