package jsonformatter

import (
	"github.com/wernsiet/morchy/internal/domain/workload"
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
