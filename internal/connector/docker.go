package connector

import (
	"context"

	containertypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type DockerClient struct {
	client *client.Client
}

func NewDockerClient() *DockerClient {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()
	return &DockerClient{client: cli}
}

func (c *DockerClient) GetContainers() ([]containertypes.Summary, error) {
	return c.client.ContainerList(context.Background(), containertypes.ListOptions{})
}
