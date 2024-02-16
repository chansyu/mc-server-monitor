package admin_console

import (
	"context"
	"log"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
)

type GCPAdminConsole struct {
	project  string
	instance string
	zone     string
}

func Open(project, instance, zone string) (*GCPAdminConsole, error) {
	return &GCPAdminConsole{
		project, instance, zone,
	}, nil
}

func (c GCPAdminConsole) Start(ctx context.Context) error {
	client, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	req := &computepb.StartInstanceRequest{
		Project:  c.project,
		Instance: c.instance,
		Zone:     c.zone,
	}
	op, err := client.Start(ctx, req)
	if err != nil {
		return err
	}

	err = op.Wait(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (c GCPAdminConsole) Restart(ctx context.Context) error {
	client, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	req := &computepb.ResetInstanceRequest{
		Project:  c.project,
		Instance: c.instance,
		Zone:     c.zone,
	}
	op, err := client.Reset(ctx, req)
	if err != nil {
		log.Fatal(err)
	}

	err = op.Wait(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func (c GCPAdminConsole) Stop(ctx context.Context) error {
	client, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	req := &computepb.StopInstanceRequest{
		Project:  c.project,
		Instance: c.instance,
		Zone:     c.zone,
	}
	op, err := client.Stop(ctx, req)
	if err != nil {
		return err
	}

	err = op.Wait(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (c GCPAdminConsole) IsOnline(ctx context.Context) (bool, error) {
	client, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return false, err
	}
	defer client.Close()

	req := &computepb.GetInstanceRequest{
		Project:  c.project,
		Instance: c.instance,
		Zone:     c.zone,
	}
	instance, err := client.Get(ctx, req)
	if err != nil {
		return false, err
	}
	return *instance.Status == "RUNNING", nil
}
