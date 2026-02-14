package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func NewControlPlaneCommand() *cobra.Command {
	cfg := &Config{}

	cmd := &cobra.Command{
		Use: "controlplane",
		RunE: func(cmd *cobra.Command, args []string) error {
			appFx := fx.New(
				fx.Provide(
					func() *Config { return cfg },
					newLogger,
					newContext,
					newDBPool,
					newWorkloadRepository,
					newUsecaseHandler,
					newRouter,
					newHTTPServer,
					newBackgroundTaskRunner,
				),
				fx.WithLogger(
					func(l *zap.Logger) fxevent.Logger {
						return &fxevent.ZapLogger{Logger: l}
					},
				),
				fx.Invoke(runServer),
				fx.Invoke(runBackgroundWorker),
			)

			startCtx, cancel := context.WithTimeout(cmd.Context(), 15*time.Second)
			defer cancel()
			if err := appFx.Start(startCtx); err != nil {
				return err
			}

			// Wait for termination signal
			quit := make(chan os.Signal, 1)
			signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
			select {
			case <-quit:
			case <-cmd.Context().Done():
			}

			stopCtx, stopCancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer stopCancel()
			return appFx.Stop(stopCtx)
		},
	}

	cmd.Flags().IntVar(&cfg.Port, "port", 8080, "HTTP server port")
	cmd.Flags().StringVar(&cfg.DBConnString, "db", os.Getenv("DATABASE_URL"), "Postgres connection string")
	cmd.Flags().IntVar(&cfg.LeaseLifetimeSec, "lease-lifetime", 30, "Lease lifetime in seconds before expiration")
	cmd.Flags().IntVar(&cfg.EventListLimit, "event-list-limit", 5, "Maximum number of events to return per workload")
	cmd.Flags().IntVar(&cfg.StuckTimeoutSec, "stuck-timeout", 90, "Timeout in seconds before marking a workload as stuck")

	return cmd
}
