package usecase

import (
	"context"

	"github.com/wernsiet/morchy/pkg/controlplane/domain/workload"
	"github.com/wernsiet/morchy/pkg/runtime"
	"go.uber.org/zap"
)

type NodeLogic interface {
	NodeDiscover(context.Context, runtime.ResourceLimits)
	NodeLaunch(context.Context, workload.Workload)
	NodeTerminate(context.Context, workload.Workload)
	NodeAck(context.Context, string)
}

type WorkloadLogic interface {
	ListWorkloads(ctx context.Context, statusEq *string, resourceLte *runtime.ResourceLimits) ([]*workload.Workload, error)
	GetWorkload(context.Context, string) (*workload.Workload, error)
	CreateWorkload(ctx context.Context, workloadSpec workload.WorkloadSpec) (*workload.Workload, error)
}

type LeaseLogic interface {
	CreateOrExtendLease(context.Context, string, string) (*workload.Lease, error)
	GetLeaseByNodeAndWorkloadID(context.Context, string, string) (*workload.Lease, error)
	ExpireLeases(context.Context) error
	DeleteLease(context.Context, string, string) error
}

type EventLogic interface {
	PushEvent(ctx context.Context, event workload.Event) error
}

type Handler interface {
	WorkloadLogic
	LeaseLogic
	EventLogic
}

type interactor struct {
	logger       *zap.Logger
	wokrloadRepo workload.Repository
}

func NewHandler(
	logger *zap.Logger,
	workloadRepo workload.Repository,
) Handler {
	return &interactor{
		logger:       logger,
		wokrloadRepo: workloadRepo,
	}
}
