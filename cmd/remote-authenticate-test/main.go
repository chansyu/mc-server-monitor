package main

import (
	"context"
	"fmt"
	"log"

	admin_console "github.com/itzsBananas/mc-server-monitor/internal/admin-console"
)

func main() {
	ctx := context.Background()

	console, err := admin_console.Open(ctx, "minecraft-626", "minecraft", "us-west2-a")
	if err != nil {
		log.Fatal(err)
	}
	defer console.Close()

	b, err := console.IsOnline()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(b)
}
