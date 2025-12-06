package usecase

import (
	"context"

	"github.com/wernsiet/morchy/internal/domain/workload"

	pkgworkload "github.com/wernsiet/morchy/pkg/workload"
)

type NodeLogic interface {
	NodeDiscover(context.Context, pkgworkload.ResourceLimits)
	NodeLaunch(context.Context, workload.Workload)
	NodeTerminate(context.Context, workload.Workload)
	NodeAck(context.Context, string)
}

type WorkloadLogic interface {
	ListWorkloads(ctx context.Context, statusEq *string, resourceLte *pkgworkload.ResourceLimits) ([]*workload.Workload, error)
	GetWorkload(context.Context, string)
	CreateWorkload(ctx context.Context, workloadSpec workload.WorkloadSpec) (*workload.Workload, error)
}

type LeaseLogic interface {
	CreateLease(context.Context, string, string) (*workload.Lease, error)
	ExtendLease(context.Context, string, string) error
	GetLeaseByNodeAndWorkloadID(context.Context, string, string) (*workload.Lease, error)
	ExpireLeases(context.Context) error
}

type Handler interface {
	WorkloadLogic
	LeaseLogic
}

type interactor struct {
	wokrloadRepo workload.Repository
}

func NewHandler(
	workloadRepo workload.Repository,
) Handler {
	return &interactor{
		wokrloadRepo: workloadRepo,
	}
}
