package jsonformatter

import (
	"time"

	"github.com/wernsiet/morchy/internal/domain/workload"
)

type LeaseResponse struct {
	WorkloadID string    `json:"workload_id" example:"some-workload-id"`
	NodeID     string    `json:"node_id" example:"some-node-id"`
	CreatedAt  time.Time `json:"created_at" example:"2025-11-03T13:45:00+03"`
	UpdatedAt  time.Time `json:"updated_at" example:"2025-11-03T13:45:00+03"`
}

func NewLeaseResponseFromDomain(l *workload.Lease) *LeaseResponse {
	return &LeaseResponse{
		NodeID:     l.NodeID,
		WorkloadID: l.WorkloadID,
		CreatedAt:  l.CreatedAt,
		UpdatedAt:  l.UpdatedAt,
	}
}
