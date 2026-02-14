package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	ginrouter "github.com/wernsiet/morchy/pkg/controlplane/implementation/gin.router"
	"github.com/wernsiet/morchy/pkg/controlplane/implementation/repository/workload"
	"github.com/wernsiet/morchy/pkg/controlplane/infrastructure"
	"github.com/wernsiet/morchy/pkg/controlplane/usecase"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func newContext(lc fx.Lifecycle) context.Context {
	ctx, cancel := context.WithCancel(context.Background())

	lc.Append(fx.Hook{
		OnStop: func(ctx2 context.Context) error {
			cancel()
			return nil
		},
	})

	return ctx
}

func newLogger() (*zap.Logger, error) {
	return zap.NewDevelopment()
}

func newDBPool(ctx context.Context, cfg *Config) (*pgxpool.Pool, error) {
	return infrastructure.NewPgxpool(ctx, cfg.DBConnString)
}

func newWorkloadRepository(dbPool *pgxpool.Pool) *workload.Repository {
	return workload.NewRepo(dbPool)
}

func newUsecaseHandler(logger *zap.Logger, workloadRepo *workload.Repository, dbPool *pgxpool.Pool, cfg *Config) usecase.Handler {
	return usecase.NewHandler(logger, workloadRepo, workload.WorkloadRepoFactory{}, dbPool, cfg.LeaseLifetimeSec, cfg.EventListLimit, cfg.StuckTimeoutSec)
}

func newRouter(logger *zap.Logger, ucHandler usecase.Handler) *gin.Engine {
	r := infrastructure.NewRouter(logger)
	rH := ginrouter.NewRouterHandler(logger, ucHandler)
	rH.SetRoutes(r)
	return r
}

func newHTTPServer(cfg *Config, r *gin.Engine) *http.Server {
	return &http.Server{
		Handler:           r,
		Addr:              fmt.Sprintf(":%d", cfg.Port),
		ReadHeaderTimeout: 5 * time.Second,
	}
}

func newBackgroundTaskRunner(logger *zap.Logger, ucHandler usecase.Handler) *infrastructure.BackgroundTaskRunner {
	bgTaskRunner := infrastructure.NewBackgroundTaskRunner(logger)
	bgTaskRunner.RegisterTask("ExpireLeases", ucHandler.ExpireLeases, 30*time.Second)
	return bgTaskRunner
}

func runServer(lc fx.Lifecycle, server *http.Server, logger *zap.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go func() {
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					logger.Error("HTTP server crashed", zap.Error(err))
				}
			}()
			logger.Info("server started", zap.String("addr", server.Addr))
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("shutting down HTTP server")
			if err := server.Shutdown(ctx); err != nil {
				logger.Error("HTTP server shutdown failed", zap.Error(err))
			}
			return nil
		},
	})
}

func runBackgroundWorker(lc fx.Lifecycle, ctx context.Context, bgRunner *infrastructure.BackgroundTaskRunner, logger *zap.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go bgRunner.Start(ctx)
			logger.Info("started background tasks")
			return nil
		},
		OnStop: func(_ context.Context) error {
			logger.Info("shutting down background tasks")
			return nil
		},
	})
}
