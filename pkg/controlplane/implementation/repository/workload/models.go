package workload

import (
	"time"

	"github.com/wernsiet/morchy/pkg/controlplane/domain/workload"
)

type dbWorkload struct {
	ID        string
	Status    string
	CreatedAt time.Time
	Spec      dbWorkloadSpec
}

type dbWorkloadSpec struct {
	ID      string
	Image   string
	CPU     uint
	RAM     uint
	Command []string
	Env     map[string]string
}

func (w *dbWorkload) ToDomain() *workload.Workload {
	return &workload.Workload{
		ID:     w.ID,
		Status: workload.WorkloadStatus(w.Status),
		Spec: workload.WorkloadSpec{
			Name:    w.ID,
			Image:   w.Spec.Image,
			CPU:     w.Spec.CPU,
			RAM:     w.Spec.RAM,
			Command: w.Spec.Command,
			Env:     w.Spec.Env,
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
