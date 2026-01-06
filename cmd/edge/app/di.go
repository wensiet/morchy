package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	ginzapcontrib "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/wernsiet/morchy/pkg/edge/implementation/controlplane"
	ginrouter "github.com/wernsiet/morchy/pkg/edge/implementation/gin.router"
	"github.com/wernsiet/morchy/pkg/edge/implementation/repository"
	"github.com/wernsiet/morchy/pkg/edge/usecase"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func newLogger() (*zap.Logger, error) {
	return zap.NewDevelopment()
}

func newRepository() *repository.Repository {
	return repository.NewRepository()
}

func newHTTPClient() *resty.Client {
	return resty.New()
}

func newControlPlaneClient(cfg *Config, http *resty.Client) *controlplane.Client {
	return controlplane.NewClient(cfg.ControlPlaneURL, http)
}

func newUsecaseHandler(logger *zap.Logger, cpClient *controlplane.Client, repo *repository.Repository) usecase.Handler {
	return usecase.NewHandler(logger, cpClient, repo)
}

func newRouter(logger *zap.Logger, ucHandler usecase.Handler) *gin.Engine {
	r := gin.New()
	r.Use(ginzapcontrib.RecoveryWithZap(logger, true))
	rH := ginrouter.NewRouterHandler(logger, ucHandler)
	rH.SetRoutes(r)

	return r
}

func newHTTPServer(r *gin.Engine) *http.Server {
	return &http.Server{
		Handler:           r,
		Addr:              fmt.Sprintf(":%d", 9999),
		ReadHeaderTimeout: 5 * time.Second,
	}
}

func runProxy(lc fx.Lifecycle, server *http.Server, logger *zap.Logger) {
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

func runEdgeSync(lc fx.Lifecycle, ucHandler usecase.Handler, logger *zap.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			logger.Info("starting edge sync loop (every 10s)")
			go func() {
				ticker := time.NewTicker(10 * time.Second)
				defer ticker.Stop()

				for range ticker.C {
					if err := ucHandler.UpsertEdges(context.Background()); err != nil {
						logger.Warn("failed to sync edges — will retry in 10s", zap.Error(err))
					} else {
						logger.Debug("edges synced successfully")
					}
				}
			}()
			return nil
		},
	})
}
