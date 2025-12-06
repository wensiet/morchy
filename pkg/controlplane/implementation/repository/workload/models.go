package workload

import (
	"encoding/json"
	"time"

	"github.com/wernsiet/morchy/pkg/controlplane/domain/workload"
	"github.com/wernsiet/morchy/pkg/runtime"
)

type dbWorkload struct {
	ID        string
	Status    string
	CreatedAt time.Time
	Container json.RawMessage
}

func (w *dbWorkload) ToDomain() *workload.Workload {
	var c runtime.Container

	if len(w.Container) > 0 {
		_ = json.Unmarshal(w.Container, &c) // TODO: handle err
	}
	return &workload.Workload{
		ID:     w.ID,
		Status: workload.WorkloadStatus(w.Status),
		Spec: workload.WorkloadSpec{
			Container: c,
		},
	}
}

type dbLease struct {
	ID         string
	NodeID     string
	WorkloadID string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (l *dbLease) ToDomain() *workload.Lease {
	return &workload.Lease{
		NodeID:     l.NodeID,
		WorkloadID: l.WorkloadID,
		CreatedAt:  l.CreatedAt,
		UpdatedAt:  l.UpdatedAt,
	}
}
