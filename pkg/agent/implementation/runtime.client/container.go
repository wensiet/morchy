package runtimeclient

import (
	"context"

	"github.com/docker/docker/api/types/container"
	"github.com/samber/oops"
	"github.com/wernsiet/morchy/pkg/agent/domain/workload"
)

func (c *RuntimeClient) CreateContainer(ctx context.Context, w workload.Workload) (string, error) {
	containerEntity, err := c.dockerClient.ContainerCreate(
		ctx,
		&container.Config{},
		nil,
		nil,
		nil,
		w.ID,
	)
	if err != nil {
		return "", oops.
			In("runtime.CreateContainer").
			With("workload_id", w.ID).
			With("action", "dockerClient.ContainerCreate").
			Wrap(err)
	}
	return containerEntity.ID, nil
}
