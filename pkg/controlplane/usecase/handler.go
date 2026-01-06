package usecase

import (
	"context"

	"github.com/wernsiet/morchy/pkg/controlplane/domain/workload"
	dbutils "github.com/wernsiet/morchy/pkg/db.utils"
	"github.com/wernsiet/morchy/pkg/runtime"
	"go.uber.org/zap"
)

type WorkloadLogic interface {
	ListWorkloads(ctx context.Context, statusEq *string, resourceLte *runtime.ResourceLimits) ([]*workload.Workload, error)
	GetWorkload(context.Context, string) (*workload.Workload, error)
	CreateWorkload(ctx context.Context, workloadSpec workload.WorkloadSpec) (*workload.Workload, error)
	DeleteWorkload(context.Context, string) error
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

type EdgeLogic interface {
	ListEdges(context.Context) ([]*workload.Edge, error)
}

type Handler interface {
	WorkloadLogic
	LeaseLogic
	EventLogic
	EdgeLogic
}

type interactor struct {
	logger            *zap.Logger
	dbPool            dbutils.DB
	wokrloadRepo      workload.Repository
	repositoryFactory workload.RepositoryFactory
}

func NewHandler(
	logger *zap.Logger,
	workloadRepo workload.Repository,
	workloadRepoFactory workload.RepositoryFactory,
	dbPool dbutils.DB,
) Handler {
	return &interactor{
		logger:            logger,
		wokrloadRepo:      workloadRepo,
		repositoryFactory: workloadRepoFactory,
		dbPool:            dbPool,
	}
}
