package usecase

import (
	"context"
	"time"

	"github.com/wernsiet/morchy/pkg/agent/domain"
	"github.com/wernsiet/morchy/pkg/agent/domain/workload"
	"github.com/wernsiet/morchy/pkg/agent/infrastructure"
	"go.uber.org/zap"
)

func (i *interactor) startWorklodLifecycle(ctx context.Context, wl workload.Workload) error {
	if err := i.ReconcileWorkload(ctx, wl); err != nil {
		return err
	}
	i.workloadSupervisor.Start(ctx, infrastructure.PeriodicTask[workload.Workload]{
		ID:       infrastructure.TaskID(wl.ID),
		Interval: 10 * time.Second,
		Input:    wl,
		Execute:  i.ReconcileWorkload,
		OnError: func(err error) {
			i.logger.Warn("workload reconciliation failed",
				zap.String(domain.SWorkloadID, wl.ID),
				zap.Error(err),
			)
		},
	})
	return nil
}

func (i *interactor) stopWorkloadLifecycle(ctx context.Context, wl workload.Workload) {
	i.workloadSupervisor.Stop(infrastructure.TaskID(wl.ID))
}
