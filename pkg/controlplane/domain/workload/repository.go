package workload

import (
	"context"

	"github.com/wernsiet/morchy/pkg/runtime"
)

type Repository interface {
	ListWorkloads(context.Context, *string, *runtime.ResourceLimits) ([]*Workload, error)
	CreateWorkload(context.Context, Workload) (*Workload, error)
	GetWorkload(context.Context, string) (*Workload, error)

	GetLease(context.Context, string, string) (*Lease, error)
	CreateLease(context.Context, string, string) (*Lease, error)
	DeleteExpiredLeases(context.Context, int) error
	DeleteLease(context.Context, string, string) error
	UpdateLease(context.Context, string, string) error
	UpsertLease(ctx context.Context, nodeID string, workloadID string) (*Lease, error)

	SaveEvent(context.Context, Event) error
}
