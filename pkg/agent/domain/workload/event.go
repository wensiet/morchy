package workload

import "github.com/wernsiet/morchy/pkg/agent/domain"

type EventPayload struct {
	WorkloadID string                   `json:"workload_id"`
	Action     string                   `json:"action"`
	Status     domain.EventActionStatus `json:"status"`
}
