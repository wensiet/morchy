package usecase

import (
	"context"

	"github.com/wernsiet/morchy/pkg/controlplane/domain"
	"github.com/wernsiet/morchy/pkg/controlplane/domain/workload"
	"go.uber.org/zap"
)

func (i *interactor) ListEdges(ctx context.Context) ([]*workload.Edge, error) {
	logger := i.logger.With(
		zap.String(domain.SDomain, domain.SWorkload),
	)

	edges, err := i.wokrloadRepo.ListEdges(ctx)
	if err != nil {
		logger.Error("failed to list edges", zap.Error(err))
		return nil, err
	}

	return edges, nil
}
