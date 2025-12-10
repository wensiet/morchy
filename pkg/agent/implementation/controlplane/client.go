package controlplane

import (
	"context"

	"github.com/go-resty/resty/v2"
	apitypes "github.com/wernsiet/morchy/pkg/controlplane/implementation/jsonformatter"
	"github.com/wernsiet/morchy/pkg/runtime"
)

type ControlPlaneClient interface {
	ListAvailableWorkloads(context.Context, runtime.ResourceLimits) ([]*apitypes.WorkloadResponse, error)
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
