package repository

import (
	"context"

	apitypes "github.com/wernsiet/morchy/pkg/controlplane/implementation/jsonformatter"
	"github.com/wernsiet/morchy/pkg/edge/domain"
)

func (r *Repository) UpsertEdges(ctx context.Context, newEdges []*apitypes.EdgeResponse) {
	// TODO: replace mutex with atomic map swap
	r.mu.Lock()
	defer r.mu.Unlock()

	r.edgeStorage = make(map[ProxyPath]*domain.Edge)

	for _, edgeResp := range newEdges {
		r.edgeStorage[ProxyPath(edgeResp.ProxyPath)] = &domain.Edge{
			ProxyPath:       edgeResp.ProxyPath,
			UpstreamAddress: edgeResp.UpstreamAddress,
		}
	}
}

func (r *Repository) FindEdge(ctx context.Context, path string) *domain.Edge {
	edge, ok := r.edgeStorage[ProxyPath(path)]
	if !ok {
		return nil
	}
	return edge
}
