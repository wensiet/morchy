package usecase

import (
	"context"
	"encoding/json"
	"os"

	"github.com/antihax/optional"
	"github.com/goccy/go-yaml"
	"github.com/jedib0t/go-pretty/v6/table"
	generated "github.com/wernsiet/morchy/pkg/mctl/generated/controlplane.go"
)

func printMultipleWorkloads(workloads []generated.JsonformatterWorkloadResponse) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"ID", "STATUS"})
	for _, wl := range workloads {
		t.AppendRows([]table.Row{
			{wl.Id, wl.Status},
		})
	}
	t.Render()
}

func (i *interactor) GetWorkloadByID(ctx context.Context, workloadID string) {
	wl, _, err := i.controlplaneClient.WorkloadsApi.ApiV1WorkloadsWorkloadIdGet(ctx, workloadID)
	if err != nil {
		panic(err)
	}
	printMultipleWorkloads([]generated.JsonformatterWorkloadResponse{wl})
}

func (i *interactor) ListWorkloads(ctx context.Context, status *string, cpu *int32, ram *int32) {
	opts := generated.WorkloadsApiApiV1WorkloadsGetOpts{}
	if status != nil {
		opts.Status = optional.NewString(*status)
	}
	if cpu != nil {
		opts.Cpu = optional.NewInt32(*cpu)
	}
	if ram != nil {
		opts.Ram = optional.NewInt32(*ram)
	}

	wl, _, err := i.controlplaneClient.WorkloadsApi.ApiV1WorkloadsGet(ctx, &opts)
	if err != nil {
		panic(err)
	}
	printMultipleWorkloads(wl)
}

func (i *interactor) CreateWorkload(
	ctx context.Context,
	raw []byte,
	isYAML bool,
) {
	var spec generated.JsonformatterWorkloadSpecRequest

	var err error
	if isYAML {
		err = yaml.Unmarshal(raw, &spec)
	} else {
		err = json.Unmarshal(raw, &spec)
	}
	if err != nil {
		panic(err)
	}

	resp, _, err := i.controlplaneClient.WorkloadsApi.
		ApiV1WorkloadsPost(ctx, spec)
	if err != nil {
		panic(err)
	}

	printMultipleWorkloads([]generated.JsonformatterWorkloadResponse{resp})
}
