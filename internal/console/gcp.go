package console

import (
	"context"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
)

type GCP struct {
	project  string
	instance string
	zone     string
}

func GCPOpen(project, instance, zone string) (*GCP, error) {
	return &GCP{
		project, instance, zone,
	}, nil
}

func (c GCP) Start(ctx context.Context) error {
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

func (c GCP) Restart(ctx context.Context) error {
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
		return err
	}

	err = op.Wait(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (c GCP) Stop(ctx context.Context) error {
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

func (c GCP) IsOnline(ctx context.Context) (bool, error) {
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

func (c GCP) Close() error {
	return nil
}
