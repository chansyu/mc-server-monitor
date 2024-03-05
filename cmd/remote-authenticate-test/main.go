package main

import (
	"context"
	"log"

	admin_console "github.com/itzsBananas/mc-server-monitor/internal/admin-console"
)

func main() {
	ctx := context.Background()

	console, err := admin_console.GCPAdminConsoleOpen("minecraft-626", "minecraft", "us-west2-a")
	if err != nil {
		log.Fatal(err)
	}

	err = console.Start(ctx)
	if err != nil {
		log.Fatal(err)
	}

}
