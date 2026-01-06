package usecase

import (
	"context"

	"github.com/wernsiet/morchy/pkg/edge/domain"
	"github.com/wernsiet/morchy/pkg/edge/implementation/controlplane"
	"go.uber.org/zap"
)

type EdgeLogic interface {
	UpsertEdges(ctx context.Context) error
	FindEdge(ctx context.Context, path string) *domain.Edge
}

type Handler interface {
	EdgeLogic
}

type interactor struct {
	logger             *zap.Logger
	controlplaneClient controlplane.ControlPlaneClient
	repository         domain.Repository
}

func NewHandler(
	logger *zap.Logger,
	controlplaneClient controlplane.ControlPlaneClient,
	repository domain.Repository,
) Handler {
	return &interactor{
		logger:             logger,
		controlplaneClient: controlplaneClient,
		repository:         repository,
	}
}
