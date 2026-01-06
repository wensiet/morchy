package app

import (
	"context"
	"os"
	"time"

	dockerclient "github.com/docker/docker/client"
	"github.com/wernsiet/morchy/pkg/agent/implementation/controlplane"
	"github.com/wernsiet/morchy/pkg/agent/implementation/repository/workload"
	"github.com/wernsiet/morchy/pkg/agent/implementation/supervisor"
	"github.com/wernsiet/morchy/pkg/agent/usecase"
	"github.com/wernsiet/morchy/pkg/runtime"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/go-resty/resty/v2"
)

func newLogger() (*zap.Logger, error) {
	return zap.NewDevelopment()
}

func newDockerRuntime() (*runtime.Client, error) {
	api, err := dockerclient.NewClientWithOpts(
		dockerclient.WithHost(os.Getenv("DOCKER_HOST")),
		dockerclient.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return nil, err
	}
	return runtime.NewClient(api), nil
}

func newHTTPClient() *resty.Client {
	return resty.New()
}

func newControlPlaneClient(cfg *Config, http *resty.Client) *controlplane.Client {
	return controlplane.NewClient(http, cfg.ControlPlaneURL, cfg.NodeID)
}

func newWorkloadRepository(cfg *Config) *workload.Repository {
	repo := workload.NewRepository()
	repo.SetResourceLimits(runtime.ResourceLimits{
		CPU: cfg.ReservedCPU,
		RAM: cfg.ReservedRAM,
	})
	return repo
}

func newWorkloadSupervisor() *supervisor.WorkloadSupervisor {
	return supervisor.NewSupervisor()
}

func newHandler(
	logger *zap.Logger,
	cp *controlplane.Client,
	rt *runtime.Client,
	repo *workload.Repository,
	wlSupervisor *supervisor.WorkloadSupervisor,
) usecase.Handler {
	return usecase.NewHandler(logger, cp, rt, repo, wlSupervisor)
}

func runLoop(lc fx.Lifecycle, logger *zap.Logger, h usecase.Handler) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				backoffLimit := 10
				for {
					err := h.LoadCurrentState(ctx)
					if err == nil {
						break
					} else {
						backoffLimit -= 1
						if backoffLimit <= 0 {
							logger.Error("loading error backoff limit reached", zap.Error(err))
							os.Exit(1)
						}
					}
				}

				if err := h.ApplyWorkloadJoin(ctx); err != nil {
					logger.Error("initial ApplyWorkloadJoin failed", zap.Error(err))
				}

				ticker := time.NewTicker(10 * time.Second)
				defer ticker.Stop()

				for {
					select {
					case <-ctx.Done():
						logger.Info("runLoop background worker stopping")
						return
					case <-ticker.C:
						if err := h.ApplyWorkloadJoin(ctx); err != nil {
							logger.Error("ApplyWorkloadJoin failed", zap.Error(err))
						}
					}
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("runLoop OnStop")
			return nil
		},
	})
}
