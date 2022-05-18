package main

import (
	"context"
	"log"
	"os"

	"github.com/bspain/funkytown/shared/redisfacade"
)

var redis_host = os.Getenv("REDIS_HOST")
var redis_port = os.Getenv("REDIS_PORT")
var ctx = context.Background()

func main() {
	log.Printf("funkytown Controller has started...")

	f := redisfacade.NewFacade(ctx, redis_host, redis_port)

	f.SetRunMetadata("a_new_run", 0)
}