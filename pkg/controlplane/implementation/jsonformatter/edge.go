package jsonformatter

import "github.com/wernsiet/morchy/pkg/controlplane/domain/workload"

type EdgeResponse struct {
	UpstreamAddress string `json:"upstream_address"`
	ProxyPath       string `json:"proxy_path"`
}

func NewEdgeFromDomain(e *workload.Edge) *EdgeResponse {
	return &EdgeResponse{
		UpstreamAddress: e.UpstreamAddress,
		ProxyPath:       e.ProxyPath,
	}
}
