package admin_console

import (
	"context"
	"log"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
)

type AdminConsole struct {
	client            *compute.InstancesClient
	context           context.Context
	requestParameters RequestParameters
}

type RequestParameters struct {
	project  string
	instance string
	zone     string
}

func Open(ctx context.Context, project, instance, zone string) (*AdminConsole, error) {
	c, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return nil, err
	}
	return &AdminConsole{
		client:  c,
		context: ctx,
		requestParameters: RequestParameters{
			project, instance, zone,
		},
	}, nil
}

func (c AdminConsole) Start() error {
	req := &computepb.StartInstanceRequest{
		Project:  c.requestParameters.project,
		Instance: c.requestParameters.instance,
		Zone:     c.requestParameters.zone,
	}
	op, err := c.client.Start(c.context, req)
	if err != nil {
		return err
	}

	err = op.Wait(c.context)
	if err != nil {
		return err
	}
	return nil
}

func (c AdminConsole) Restart() error {
	req := &computepb.ResetInstanceRequest{
		Project:  c.requestParameters.project,
		Instance: c.requestParameters.instance,
		Zone:     c.requestParameters.zone,
	}
	op, err := c.client.Reset(c.context, req)
	if err != nil {
		log.Fatal(err)
	}

	err = op.Wait(c.context)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func (c AdminConsole) Stop() error {
	req := &computepb.StopInstanceRequest{
		Project:  c.requestParameters.project,
		Instance: c.requestParameters.instance,
		Zone:     c.requestParameters.zone,
	}
	op, err := c.client.Stop(c.context, req)
	if err != nil {
		return err
	}

	err = op.Wait(c.context)
	if err != nil {
		return err
	}
	return nil
}

func (c AdminConsole) IsOnline() (bool, error) {
	req := &computepb.GetInstanceRequest{
		Project:  c.requestParameters.project,
		Instance: c.requestParameters.instance,
		Zone:     c.requestParameters.zone,
	}
	instance, err := c.client.Get(c.context, req)
	if err != nil {
		return false, err
	}
	return *instance.Status == "RUNNING", nil
}

func (c AdminConsole) Close() {
	c.Close()
}
