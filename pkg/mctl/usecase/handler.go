package usecase

import (
	"context"

	generated "github.com/wernsiet/morchy/pkg/mctl/generated/controlplane.go"
)

type WorkloadLogic interface {
	GetWorkloadByID(ctx context.Context, workloadID string)
	ListWorkloads(ctx context.Context, status *string, cpu *int32, ram *int32)
	CreateWorkload(ctx context.Context, raw []byte, isYAML bool)
	DeleteWorkload(ctx context.Context, workloadID string)
}

type Handler interface {
	WorkloadLogic
}

type interactor struct {
	controlplaneClient *generated.APIClient
}

func NewHandler(
	controlplaneClient *generated.APIClient,
) Handler {
	return &interactor{
		controlplaneClient: controlplaneClient,
	}
}
