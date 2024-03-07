package admin_console

import (
	"context"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type LocalAdminConsole struct {
	client      *client.Client
	containerID string
}

// start, restart, stop, maybe see if it's on too?
// only to be used locally...
func LocalAdminConsoleOpen(containerID string) (*LocalAdminConsole, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	return &LocalAdminConsole{
		client:      cli,
		containerID: containerID,
	}, nil
}

func (c *LocalAdminConsole) Start(ctx context.Context) error {
	return c.client.ContainerStart(ctx, c.containerID, container.StartOptions{})
}

func (c *LocalAdminConsole) Restart(ctx context.Context) error {
	return c.client.ContainerRestart(ctx, c.containerID, container.StopOptions{})
}

func (c *LocalAdminConsole) Stop(ctx context.Context) error {
	return c.client.ContainerStop(ctx, c.containerID, container.StopOptions{})
}

func (c *LocalAdminConsole) IsOnline(ctx context.Context) (bool, error) {
	stats, err := c.client.ContainerInspect(ctx, c.containerID)
	if err != nil {
		return false, nil
	}
	return stats.ContainerJSONBase.State.Running, nil
}

func (c *LocalAdminConsole) Close() error {
	return c.client.Close()
}
