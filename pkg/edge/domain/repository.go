package domain

import (
	"context"

	apitypes "github.com/wernsiet/morchy/pkg/controlplane/implementation/jsonformatter"
)

type Repository interface {
	UpsertEdges(ctx context.Context, newEdges []*apitypes.EdgeResponse)
	FindEdge(ctx context.Context, path string) *Edge
}
