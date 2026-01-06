package workload

import (
	"database/sql"
	"encoding/json"
	"fmt"
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
	ID            string
	Image         string
	CPU           uint
	RAM           uint
	Command       []string
	Env           map[string]string
	ContainerPort sql.NullInt32
	HostPort      sql.NullInt32
}

func (w *dbWorkload) ToDomain() *workload.Workload {
	wl := workload.Workload{
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
	if w.Spec.ContainerPort.Valid {
		containerPort := int(w.Spec.ContainerPort.Int32)
		wl.Spec.ContainerPort = &containerPort
	}
	if w.Spec.HostPort.Valid {
		hostPort := int(w.Spec.HostPort.Int32)
		wl.Spec.HostPort = &hostPort
	}
	return &wl
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

type dbEvent struct {
	ID         string
	SourceID   string
	NodeID     string
	Payload    json.RawMessage
	ProducedAt time.Time
	CreatedAt  time.Time
}

func (e *dbEvent) ToDomain() *workload.Event {
	return &workload.Event{
		ID:         e.ID,
		SourceID:   e.SourceID,
		NodeID:     e.NodeID,
		Payload:    e.Payload,
		ProducedAt: e.ProducedAt,
		CreatedAt:  e.CreatedAt,
	}
}

type dbEdge struct {
	WorkloadID string
	NodeID     string
	HostPort   int
}

func (e *dbEdge) ToDomain() *workload.Edge {
	return &workload.Edge{
		UpstreamAddress: fmt.Sprintf("localhost:%d", e.HostPort),
		ProxyPath:       fmt.Sprintf("/%s", e.WorkloadID),
	}
}
