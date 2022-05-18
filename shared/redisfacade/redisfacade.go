package redisfacade

import (
	"context"
	"log"

	"github.com/bspain/funkytown/shared/model"
	"github.com/go-redis/redis/v8"
)

// runmetadata is a Redis Hash (object) that holds metadata about the run
const runmetadata = "runmeta"

type RedisFacade struct {
	context context.Context
	client  *redis.Client
}

func NewFacade(context context.Context, host string, port string) RedisFacade {
	r := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: "",
		DB:       0,
	})

	f := RedisFacade{
		context: context,
		client:  r,
	}

	return f
}

func (f RedisFacade) SetRunMetadata(runid string, commandcount int) {
	h := make(map[string]interface{})

	h[string(model.KeyRunMetaId)] = runid
	h[string(model.KeyRunMetaTasksRemaining)] = commandcount
	h[string(model.KeyRunMetaTasksFinished)] = 0
	h[string(model.KeyRunMetaFinished)] = false

	_, err := f.client.HMSet(f.context, runmetadata, h).Result()
	if err != nil {
		log.Fatalf("Unable to create hash for %v, %v", runmetadata, err)
	}

	log.Printf("SetRunMetadata: run metadata set successfully.")
}
