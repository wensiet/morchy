package runtimeclient

import "github.com/docker/docker/client"

type RuntimeClient struct {
	dockerClient *client.Client
}

func NewRuntimeClient(dockerClient *client.Client) *RuntimeClient {
	return &RuntimeClient{
		dockerClient: dockerClient,
	}
}
