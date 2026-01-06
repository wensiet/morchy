package runtime

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

type RuntimeClient interface {
	CreateContainer(context.Context, Container) (string, error)
	StartContainer(context.Context, string) error
	StopContainer(context.Context, string) error
	RemoveContainer(context.Context, string) error
	ListContainers(context.Context, *ContainerFilters) ([]*ContainerBrief, error)
	GetContainerStatus(context.Context, string) (string, error)
}

type Client struct {
	dockerAPI   *client.Client
	stopTimeout int
}

func NewClient(dockerAPI *client.Client) *Client {
	return &Client{
		dockerAPI:   dockerAPI,
		stopTimeout: 1,
	}
}

func (c *Client) CreateContainer(ctx context.Context, cnt Container) (string, error) {
	var portBindings nat.PortMap
	if cnt.NetConfig != nil {
		containerPort := nat.Port(fmt.Sprintf("%d/tcp", cnt.NetConfig.ContainerPort))
		portBindings = nat.PortMap{
			containerPort: []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: fmt.Sprintf("%d", cnt.NetConfig.HostPort),
				},
			},
		}
	}

	resp, err := c.dockerAPI.ContainerCreate(
		ctx,
		&container.Config{
			Image:  cnt.Image,
			Cmd:    cnt.Command,
			Env:    buildEnv(cnt.Env),
			Labels: cnt.Labels,
		},
		&container.HostConfig{
			Resources: container.Resources{
				NanoCPUs: int64(cnt.Resources.CPU) * 1000 * 1000, // 1 millicore = 1000 * 1000 * 1000 nanocores
				Memory:   int64(cnt.Resources.RAM) * 1024 * 1024, // 1 MB = 1024 * 1024 bytes
			},
			PortBindings: portBindings,
		},
		nil,
		nil,
		cnt.Name,
	)
	if err != nil {
		return "", err
	}

	return resp.ID, nil
}

func (c *Client) StartContainer(ctx context.Context, id string) error {
	return c.dockerAPI.ContainerStart(ctx, id, container.StartOptions{})
}

func buildEnv(envs []EnvVar) []string {
	out := make([]string, 0, len(envs))
	for _, e := range envs {
		out = append(out, e.Name+"="+e.Value)
	}
	return out
}

func (c *Client) StopContainer(ctx context.Context, id string) error {
	return c.dockerAPI.ContainerStop(ctx, id, container.StopOptions{Timeout: &c.stopTimeout})
}

func (c *Client) RemoveContainer(ctx context.Context, id string) error {
	return c.dockerAPI.ContainerRemove(ctx, id, container.RemoveOptions{})
}

func (c *Client) GetContainerStatus(ctx context.Context, id string) (string, error) {
	containerInspect, err := c.dockerAPI.ContainerInspect(ctx, id)
	if err != nil {
		return "", err
	}
	return containerInspect.State.Status, nil
}

func (c *Client) ListContainers(ctx context.Context, listFitlers *ContainerFilters) ([]*ContainerBrief, error) {
	dockerFilters := filters.NewArgs()
	if listFitlers != nil {
		for key, value := range listFitlers.Labels {
			dockerFilters.Add("label", key+"="+value)
		}
	}

	containers, err := c.dockerAPI.ContainerList(
		ctx,
		container.ListOptions{
			Filters: dockerFilters,
			All:     true,
		},
	)
	if err != nil {
		return nil, err
	}

	containerEntities := make([]*ContainerBrief, 0, len(containers))
	for _, cnt := range containers {
		containerEntities = append(
			containerEntities,
			&ContainerBrief{
				Name:   cnt.Names[0],
				Image:  cnt.Image,
				Labels: cnt.Labels,
			},
		)
	}
	return containerEntities, nil
}
