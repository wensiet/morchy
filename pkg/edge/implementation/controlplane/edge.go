package controlplane

import (
	"context"
	"fmt"
	"net/http"

	apitypes "github.com/wernsiet/morchy/pkg/controlplane/implementation/jsonformatter"
)

func (c *Client) ListEdges(ctx context.Context) ([]*apitypes.EdgeResponse, error) {
	var edges []*apitypes.EdgeResponse

	request, err := c.httpClient.R().
		SetContext(ctx).
		SetHeader("Accept", "application/json").
		SetResult(&edges).
		Get(c.baseURL + "/api/v1/edges")
	if err != nil {
		return nil, err
	}

	if request.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", request.StatusCode())
	}

	return edges, nil
}
