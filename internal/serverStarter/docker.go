package serverStarter

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// For use when the Minecraft server and the web server is on the same machine
type DockerClient struct {
	client      *client.Client
	context     context.Context
	containerID string
}

// start, restart, stop, maybe see if it's on too?
// only to be used locally...
func DockerOpen(containerID string) (*DockerClient, error) {
	cli, err := client.NewClientWithOpts()
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	return &DockerClient{
		client:      cli,
		context:     ctx,
		containerID: containerID,
	}, nil
}

func (c DockerClient) Start() error {
	return c.client.ContainerStart(c.context, c.containerID, types.ContainerStartOptions{})
}

func (c DockerClient) Stop() error {
	return c.client.ContainerStop(c.context, c.containerID, container.StopOptions{})
}

func (c DockerClient) Restart() error {
	return c.client.ContainerRestart(c.context, c.containerID, container.StopOptions{})
}

// https://pkg.go.dev/github.com/docker/docker@v24.0.7+incompatible/api/types#ContainerJSON
func (c DockerClient) Ready() bool {
	stats, err := c.client.ContainerInspect(c.context, c.containerID)
	if err != nil {
		return false
	}
	fmt.Println(stats.State)
	return true
}

func (c DockerClient) Close() error {
	return c.client.Close()
}
