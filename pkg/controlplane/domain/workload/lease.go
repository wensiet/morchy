package workload

import "time"

type Lease struct {
	NodeID     string
	WorkloadID string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
