package main

import (
	"context"
	"fmt"
	"log"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
)

func main() {
	ctx := context.Background()
	// This snippet has been automatically generated and should be regarded as a code template only.
	// It will require modifications to work:
	// - It may require correct/in-range values for request initialization.
	// - It may require specifying regional endpoints when creating the service client as shown in:
	//   https://pkg.go.dev/cloud.google.com/go#hdr-Client_Options
	c, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	req := &computepb.StartInstanceRequest{
		Project:  "minecraft-626",
		Instance: "minecraft",
		Zone:     "us-west2-a",
		// TODO: Fill request struct fields.
		// See https://pkg.go.dev/cloud.google.com/go/compute/apiv1/computepb#StartInstanceRequest.
	}
	op, err := c.Start(ctx, req)
	if err != nil {
		log.Fatal(err)
	}

	err = op.Wait(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Success")
}
