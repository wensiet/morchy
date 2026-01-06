package usecase

import (
	"context"

	"github.com/wernsiet/morchy/pkg/edge/domain"
	"go.uber.org/zap"
)

func (i *interactor) UpsertEdges(ctx context.Context) error {
	edges, err := i.controlplaneClient.ListEdges(ctx)
	if err != nil {
		i.logger.Error("failed to list edges", zap.Error(err))
		return err
	}

	i.repository.UpsertEdges(ctx, edges)

	return nil
}

func (i *interactor) FindEdge(ctx context.Context, path string) *domain.Edge {
	return i.repository.FindEdge(ctx, path)
}
