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

func NewEdgeCommand() *cobra.Command {
	cfg := &Config{}

	cmd := &cobra.Command{
		Use: "edge",
		RunE: func(cmd *cobra.Command, args []string) error {
			appFx := fx.New(
				fx.Provide(
					func() *Config { return cfg },
					newLogger,
					newRepository,
					newHTTPClient,
					newControlPlaneClient,
					newUsecaseHandler,
					newRouter,
					newHTTPServer,
				),
				fx.WithLogger(
					func(l *zap.Logger) fxevent.Logger {
						return &fxevent.ZapLogger{Logger: l}
					},
				),
				fx.Invoke(
					runProxy,
					runEdgeSync,
				),
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

	cmd.Flags().StringVar(&cfg.ControlPlaneURL, "controlplane", "", "ControlPlane URL")

	return cmd
}
