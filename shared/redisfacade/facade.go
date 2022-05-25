package redisfacade

import (
	"context"
	"fmt"
	"log"

	"github.com/bspain/funkytown/shared/model"
	"github.com/go-redis/redis/v8"
)

// RUN_METADATA is a Redis Hash (object) that holds metadata about the run
const RUN_METADATA = "runmeta"
const TASK_QUEUE = "taskqueue"

const TASKSTATUS_WAITING = "waiting"
const TASKSTATUS_PROCESSING = "processing"
const TASKSTATUS_COMPLETE = "complete"

const TASKRESULT_UNKNOWN = "unknown"

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

// SetRunMetadata will set (HMSET) the runmeta hash object
func (f RedisFacade) SetRunMetadata(runid string, taskcount int) {
	keys := GetRunMetadataKeys()
	h := make(map[string]interface{})

	h[string(keys.RunID)] = runid
	h[string(keys.TasksRemaining)] = taskcount
	h[string(keys.TasksFinished)] = 0
	h[string(keys.Finished)] = false

	_, err := f.client.HMSet(f.context, RUN_METADATA, h).Result()
	if err != nil {
		log.Fatalf("Unable to create hash for %v, %v", RUN_METADATA, err)
	}

	log.Printf("SetRunMetadata: run metadata set successfully.")
}

// SetTaskMetadata will set (HMSET) the meta hash object for a task and return the hash key 
func (f RedisFacade) SetTaskMetadata(groupname string, taskindex int, task model.Task) string {
	keys := GetTaskMetadataKeys()
	h := make(map[string]interface{})

	h[string(keys.Group)] = groupname
	h[string(keys.Spec)] = task.Spec
	h[string(keys.Viewport)] = task.Viewport
	h[string(keys.Browser)] = task.Browser
	h[string(keys.Status)] = TASKSTATUS_WAITING
	h[string(keys.Result)] = TASKRESULT_UNKNOWN
	h[string(keys.StartTime)] = 0
	h[string(keys.Iterations)] = 0
	h[string(keys.Duration)] = 0
	h[string(keys.DurationString)] = "0s"

	key := fmt.Sprintf("task:%v:%v", groupname, taskindex)
	_, err := f.client.HMSet(f.context, key, h).Result()
	if err != nil {
		log.Fatalf("Unable to create task metadata for %v, %v", key, err)
	}

	log.Printf("SetTaskMetadata: task metadata for %v set successfully.", key)

	return key
}

// PushTask will push (LPUSH) a task key to the workqueue
func (f RedisFacade) PushTask(key string) {
	_, err := f.client.LPush(f.context, TASK_QUEUE, key).Result()
	if err != nil {
		log.Fatalf("unable to push to taskqueue: %v", err)
	}

	log.Printf("PushTask: task %v pushed successfully.", key)
}

// GetTaskStatus will return a tasks status (TASKSTATUS_WAITING, TASKSTATUS_PROCESSING, or TASKSTATUS_COMPLETE )
func (f RedisFacade) GetTaskStatus(key string) string {
	keys := GetTaskMetadataKeys()
	res, err := f.client.HGet(f.context, key, string(keys.Status)).Result()
	if err != nil {
		log.Fatalf("Unable to get task status for: %v, %v", key, err)
	}

	return res
}
