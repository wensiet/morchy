package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	ginrouter "github.com/wernsiet/morchy/internal/implementation/gin.router"
	"github.com/wernsiet/morchy/internal/implementation/repository/workload"
	"github.com/wernsiet/morchy/internal/infrastructure"
	"github.com/wernsiet/morchy/internal/usecase"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbPool, err := infrastructure.NewPgxpool(ctx)
	if err != nil {
		panic(err)
	}

	ucHandler := usecase.NewHandler(workload.NewRepo(dbPool))

	bgTaskRunner := infrastructure.NewBackgroundTaskRunner(logger)
	bgTaskRunner.RegisterTask("ExpireLeases", ucHandler.ExpireLeases, 30*time.Second)
	go bgTaskRunner.Start(ctx)

	r := infrastructure.NewRouter(logger)
	rH := ginrouter.NewRouterHandler(logger, ucHandler)
	rH.SetRoutes(r)

	server := http.Server{
		Handler:           r,
		Addr:              ":8080",
		ReadHeaderTimeout: 5 * time.Second,
	}
	httpServerWaitCh := make(chan error, 1)

	go func() {
		if err := server.ListenAndServe(); err != nil {
			httpServerWaitCh <- err
			close(httpServerWaitCh)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	select {
	case <-quit:
		logger.Info("graceful shutting down...")
	case err := <-httpServerWaitCh:
		logger.Error("instant shutdown on critical error", zap.Error(err))
	}

	cancel()
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("failed to gracefuly stop HTTP server", zap.Error(err))
	}

	logger.Info("shutdown complete")
}
