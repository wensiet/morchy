package controlplane

import (
	"context"

	"github.com/go-resty/resty/v2"
	"github.com/wernsiet/morchy/pkg/agent/domain/workload"
	"github.com/wernsiet/morchy/pkg/runtime"
)

type ControlPlaneClient interface {
	ListAvailableWorkloads(context.Context, runtime.ResourceLimits) ([]*workload.Workload, error)
	CreateWorkloadLease(context.Context, string) error
	ExtendWorkloadLease(context.Context, string) error
}

type Client struct {
	baseURL    string
	httpClient *resty.Client
}

func NewClient(httpClient *resty.Client, baesURL string) *Client {
	return &Client{
		httpClient: httpClient,
		baseURL:    baesURL,
	}
}
