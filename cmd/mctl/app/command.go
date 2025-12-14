package app

import (
	"context"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	swagger "github.com/wernsiet/morchy/pkg/mctl/generated/controlplane.go"
	"github.com/wernsiet/morchy/pkg/mctl/usecase"
)

func NewMCTLCommand() *cobra.Command {
	controlPlaneURL := os.Getenv("CONTROL_PLANE_URL")
	if controlPlaneURL == "" {
		controlPlaneURL = "http://localhost:8080"
	}

	handler := usecase.NewHandler(
		swagger.NewAPIClient(
			&swagger.Configuration{
				BasePath: controlPlaneURL,
			},
		),
	)

	cmd := &cobra.Command{
		Use:   "mctl",
		Short: "Morchy CLI",
	}

	getCmd := &cobra.Command{
		Use:   "get",
		Short: "Get resources",
	}
	getCmd.AddCommand(
		newGetWorkloadCommand(handler),
		newListWorkloadsCommand(handler),
	)

	applyCmd := &cobra.Command{
		Use:   "apply",
		Short: "Apply resource",
	}
	applyCmd.AddCommand(
		newCreateWorkloadCommand(handler),
	)

	cmd.AddCommand(getCmd)
	cmd.AddCommand(applyCmd)

	return cmd
}

func newGetWorkloadCommand(handler usecase.Handler) *cobra.Command {
	return &cobra.Command{
		Use:   "workload [id]",
		Short: "Get a workload by ID",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			workloadID := args[0]
			handler.GetWorkloadByID(ctx, workloadID)
		},
	}
}

func newListWorkloadsCommand(handler usecase.Handler) *cobra.Command {
	listCmd := &cobra.Command{
		Use:   "workloads",
		Short: "List workloads",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()

			// read optional flags
			status, _ := cmd.Flags().GetString("status")
			cpuStr, _ := cmd.Flags().GetString("cpu")
			ramStr, _ := cmd.Flags().GetString("ram")

			var cpuPtr *int32
			var ramPtr *int32

			if cpuStr != "" {
				c, err := strconv.ParseInt(cpuStr, 10, 32)
				if err == nil {
					c32 := int32(c)
					cpuPtr = &c32
				}
			}
			if ramStr != "" {
				r, err := strconv.ParseInt(ramStr, 10, 32)
				if err == nil {
					r32 := int32(r)
					ramPtr = &r32
				}
			}

			var statusPtr *string
			if status != "" {
				statusPtr = &status
			}

			handler.ListWorkloads(ctx, statusPtr, cpuPtr, ramPtr)
		},
	}

	listCmd.Flags().String("status", "", "Filter by status")
	listCmd.Flags().String("cpu", "", "Filter by CPU count")
	listCmd.Flags().String("ram", "", "Filter by RAM size")

	return listCmd
}

func newCreateWorkloadCommand(handler usecase.Handler) *cobra.Command {
	var file string

	cmd := &cobra.Command{
		Use:   "workload",
		Short: "Create a workload",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()

			var data []byte
			var err error

			if file != "" {
				data, err = os.ReadFile(file)
				if err != nil {
					panic(err)
				}
			} else {
				data, err = io.ReadAll(os.Stdin)
				if err != nil {
					panic(err)
				}
				if len(data) == 0 {
					panic("no input provided (use -f or stdin)")
				}
			}

			isYAML := strings.HasSuffix(file, ".yaml") || strings.HasSuffix(file, ".yml")
			handler.CreateWorkload(ctx, data, isYAML)
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", "", "Workload spec file (yaml or json)")

	return cmd

}
