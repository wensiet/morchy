package usecase

import (
	"context"

	"github.com/wernsiet/morchy/pkg/agent/domain/workload"
	"github.com/wernsiet/morchy/pkg/agent/implementation/controlplane"
	"github.com/wernsiet/morchy/pkg/runtime"
	"go.uber.org/zap"
)

type JoinLogic interface {
	ApplyWorkloadJoin(context.Context) error
}

type Handler interface {
	JoinLogic
}

type interactor struct {
	logger             *zap.Logger
	controlplaneClient controlplane.ControlPlaneClient
	runtimeClient      runtime.RuntimeClient
	workloadRepo       workload.Repository
}

func NewHandler(
	logger *zap.Logger,
	controlplaneClient controlplane.ControlPlaneClient,
	runtimeClient runtime.RuntimeClient,
	workloadRepo workload.Repository,
) Handler {
	return &interactor{
		logger:             logger,
		controlplaneClient: controlplaneClient,
		runtimeClient:      runtimeClient,
		workloadRepo:       workloadRepo,
	}
}
