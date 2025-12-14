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

func NewAgentCommand() *cobra.Command {
	cfg := &Config{}

	root := &cobra.Command{
		Use: "agent",
		RunE: func(cmd *cobra.Command, args []string) error {
			app := fx.New(
				fx.Provide(
					func() *Config { return cfg },
					newLogger,
					newDockerRuntime,
					newHTTPClient,
					newControlPlaneClient,
					newWorkloadRepository,
					newWorkloadSupervisor,
					newHandler,
				),
				fx.WithLogger(
					func(l *zap.Logger) fxevent.Logger {
						return &fxevent.ZapLogger{Logger: l}
					},
				),
				fx.Invoke(runLoop),
			)

			if err := app.Start(cmd.Context()); err != nil {
				return err
			}

			sigCh := make(chan os.Signal, 1)
			signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
			select {
			case <-cmd.Context().Done():
			case <-sigCh:
			}

			stopCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			return app.Stop(stopCtx)
		},
	}
	root.Flags().StringVar(&cfg.NodeID, "node-id", "", "Node ID")
	root.Flags().StringVar(&cfg.ControlPlaneURL, "controlplane", "", "ControlPlane URL")
	root.Flags().UintVar(&cfg.ReservedRAM, "reserved-ram", 0, "Reserved RAM bytes")
	root.Flags().UintVar(&cfg.ReservedCPU, "reserved-cpu", 0, "Reserved CPU units")

	return root
}
