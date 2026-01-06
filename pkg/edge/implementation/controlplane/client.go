package controlplane

import (
	"context"

	"github.com/go-resty/resty/v2"
	apitypes "github.com/wernsiet/morchy/pkg/controlplane/implementation/jsonformatter"
)

type ControlPlaneClient interface {
	ListEdges(context.Context) ([]*apitypes.EdgeResponse, error)
}

type Client struct {
	baseURL    string
	httpClient *resty.Client
}

func NewClient(
	baseURL string,
	httpClient *resty.Client,
) *Client {
	return &Client{
		baseURL:    baseURL,
		httpClient: httpClient,
	}
}
