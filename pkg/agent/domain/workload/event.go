package workload

import "github.com/wernsiet/morchy/pkg/agent/domain"

type EventPayload struct {
	Action string                   `json:"action"`
	Status domain.EventActionStatus `json:"status"`
}
