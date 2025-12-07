package runtime

import (
	"context"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type RuntimeClient interface {
	CreateContainer(context.Context, Container) (string, error)
	StartContainer(context.Context, string) error
	StopContainer(context.Context, string) error
	RemoveContainer(context.Context, string) error
	GetContainerStatus(context.Context, string) (string, error)
}

type Client struct {
	dockerAPI *client.Client
}

func NewClient(dockerAPI *client.Client) *Client {
	return &Client{
		dockerAPI: dockerAPI,
	}
}

func (c *Client) CreateContainer(ctx context.Context, cnt Container) (string, error) {
	resp, err := c.dockerAPI.ContainerCreate(
		ctx,
		&container.Config{
			Image: cnt.Image,
			Cmd:   cnt.Command,
			Env:   buildEnv(cnt.Env),
		},
		&container.HostConfig{
			Resources: container.Resources{
				NanoCPUs: int64(cnt.Resources.CPU) * 1000 * 1000, // 1 millicore = 1000 * 1000 * 1000 nanocores
				Memory:   int64(cnt.Resources.RAM) * 1024 * 1024, // 1 MB = 1024 * 1024 bytes
			},
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
	return c.dockerAPI.ContainerStop(ctx, id, container.StopOptions{})
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
