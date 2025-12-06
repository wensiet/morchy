package workload

import (
	"context"

	pkgworkload "github.com/wernsiet/morchy/pkg/workload"
)

type Repository interface {
	ListWorkloads(context.Context, *string, *pkgworkload.ResourceLimits) ([]*Workload, error)
	CreateWorklod(context.Context, Workload) (*Workload, error)
	GetWorkload(context.Context, string) (*Workload, error)

	GetLease(context.Context, string, string) (*Lease, error)
	CreateLease(context.Context, string, string) (*Lease, error)
	DeleteExpiredLeases(context.Context, int) error
	UpdateLease(context.Context, string, string) error
}
