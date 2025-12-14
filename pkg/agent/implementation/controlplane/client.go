package controlplane

import (
	"context"

	"github.com/go-resty/resty/v2"
	apitypes "github.com/wernsiet/morchy/pkg/controlplane/implementation/jsonformatter"
	"github.com/wernsiet/morchy/pkg/runtime"
)

type ControlPlaneClient interface {
	ListAvailableWorkloads(context.Context, runtime.ResourceLimits) ([]*apitypes.WorkloadResponse, error)
	CreateOrExtendWorkloadLease(context.Context, string) error
	DeleteWorkloadLease(context.Context, string) error
	PushEvent(context.Context, apitypes.EventCreateRequest) error
}

type Client struct {
	baseURL    string
	httpClient *resty.Client
	nodeID     string
}

func NewClient(httpClient *resty.Client, baseURL string, nodeID string) *Client {
	return &Client{
		httpClient: httpClient,
		baseURL:    baseURL,
		nodeID:     nodeID,
	}
}
