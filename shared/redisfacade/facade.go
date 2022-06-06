package redisfacade

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

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
const TASKRESULT_PASSED = "passed"
const TASKRESULT_FAILED = "failed"
const TASKRESULT_ERROR = "error"

// ----------------------
type RedisFacadeErrorType int

const (
	None RedisFacadeErrorType = iota
	TcpHostErr
	TcpCxErr
	RedisNil
)
// ----------------------

// RedisFacade represents a REDIS DB connection and supporting functions for que/dequ workqueue tasks.
type RedisFacade struct {
	context context.Context
	client  *redis.Client
}

// NewFacade will create a new facade to the REDIS DB
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

// GetRunMetadata will return up to date information about the run
func (f RedisFacade) GetRunMetadata() (model.RunMetadata, error, RedisFacadeErrorType) {
	keys := GetRunMetadataKeys()
	obj, err := f.client.HGetAll(f.context, RUN_METADATA).Result()

	if err != nil {
		if strings.Contains(err.Error(), "no such host") {
			return model.RunMetadata{}, err, TcpHostErr
		}
		if strings.Contains(err.Error(), "connection refused") {
			return model.RunMetadata{}, err, TcpCxErr
		}
		log.Fatalf("Unable to get run metadata %v", err)
	}

	if len(obj) == 0 {
		return model.RunMetadata{}, fmt.Errorf("run metadata empty"), RedisNil
	}

	t_remaining, err := strconv.Atoi(obj[string(keys.TasksRemaining)])
	if err != nil {
		log.Fatalf("unable to parse tasksremaining from run metadata: %v", err)
	}

	t_finished, err := strconv.Atoi(obj[string(keys.TasksFinished)])
	if err != nil {
		log.Fatalf("unable to parse tasksfinished from run metadata: %v", err)
	}

	finished, err := strconv.ParseBool(obj[string(keys.Finished)])
	if err != nil {
		log.Fatalf("unable to parse finished from run metadata: %v", err)
	}

	return model.RunMetadata{
		Key: RUN_METADATA,
		Runid: obj[string(keys.RunID)],
		TasksRemaining: t_remaining,
		TasksFinished:  t_finished,
		Finished: finished,
	}, nil, None
}

// UpdateRunMetadataTaskCount will set the number of tasks accordingly.  This function will also set the "finished" flag accordingly, once all tasks are noted as finished.
func (f RedisFacade) UpdateRunMetadataTaskCount(remaining int, finished int) {

	keys := GetRunMetadataKeys()

	runfinished := false
	if (remaining == 0) {
		runfinished = true
	} 

	values := []string{
		string(keys.TasksRemaining),
		fmt.Sprint(remaining),
		string(keys.TasksFinished),
		fmt.Sprint(finished),
		string(keys.Finished),
		fmt.Sprint(runfinished),
	}

	_, err := f.client.HSet(f.context, RUN_METADATA, values).Result()
	if err != nil {
		log.Fatalf("unable to update run metadata task count. %v", err)
	}

	log.Printf("updateRunMetadataTaskCount: remaining %v, finished %v, run finished %v", remaining, finished, runfinished)
}

// SetTaskMetadata will set (HMSET) the meta hash object for a task and return the hash key 
func (f RedisFacade) SetTaskMetadata(key string, task model.Task) {
	keys := GetTaskMetadataKeys()
	h := make(map[string]interface{})

	h[string(keys.Key)] = key
	h[string(keys.Spec)] = task.Spec
	h[string(keys.Viewport)] = task.Viewport
	h[string(keys.Browser)] = task.Browser
	h[string(keys.Status)] = TASKSTATUS_WAITING
	h[string(keys.Result)] = TASKRESULT_UNKNOWN
	h[string(keys.StartTime)] = 0
	h[string(keys.Duration)] = 0
	h[string(keys.DurationString)] = "0s"

	_, err := f.client.HMSet(f.context, key, h).Result()
	if err != nil {
		log.Fatalf("Unable to create task metadata for %v, %v", key, err)
	}

	log.Printf("SetTaskMetadata: task metadata for %v set successfully.", key)
}


// PushTask will push (LPUSH) a task key to the queue
func (f RedisFacade) PushTask(key string) {
	_, err := f.client.LPush(f.context, TASK_QUEUE, key).Result()
	if err != nil {
		log.Fatalf("unable to push to taskqueue: %v", err)
	}

	log.Printf("PushTask: task %v pushed successfully.", key)
}

// PopTask will pop (RPOP) a task key from the queue
func (f RedisFacade) PopTask() (empty bool, key string) {
	key, err := f.client.RPop(f.context, TASK_QUEUE).Result()
	if err == redis.Nil {
		// queue is empty
		return true, ""
	} else if err != nil {
		log.Fatalf("unable to get task from queue: %v", err)
	}

	log.Printf("PopTask: task %v popd successfully.", key)

	return false, key
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

// GetTaskMetadata will return all tasks metadata for a task
func (f RedisFacade) GetTaskMetadata(key string) model.TaskMetadata {
	keys := GetTaskMetadataKeys()
	obj, err := f.client.HGetAll(f.context, key).Result()
	if err != nil {
		log.Fatalf("unable to get task metadata for %v: %v", key, err)
	}

	start_time, err := strconv.ParseInt(obj[string(keys.StartTime)], 10, 64)
	if err != nil {
		log.Fatalf("unable to parse starttime from task %v: %v", key, err)
	}

	duration, err := strconv.ParseInt(obj[string(keys.Duration)], 10, 64)
	if err != nil {
		log.Fatalf("unable to parse duration from task %v: %v", key, err)
	}

	return model.TaskMetadata{
		Key: obj[string(keys.Key)],
		Spec: obj[string(keys.Spec)],
		Viewport: obj[string(keys.Viewport)],
		Browser: obj[string(keys.Browser)],
		Status: obj[string(keys.Status)],
		Result: obj[string(keys.Result)],
		StartTime: time.Unix(start_time, 0),
		Duration: duration,
		DurationString: obj[string(keys.DurationString)],
	}
}

// SetTaskAsProcessing will denote the current time and mark the task as TASKSTATUS_PROCESSING
func (f RedisFacade) SetTaskAsProcessing(task *model.TaskMetadata) {
	keys := GetTaskMetadataKeys()
	task.Status = TASKSTATUS_PROCESSING
	task.StartTime = time.Now()

	values := []string{
		string(keys.Status),
		string(task.Status),
		string(keys.StartTime),
		fmt.Sprint(task.StartTime.Unix()),
	}

	_, err := f.client.HSet(f.context, task.Key, values).Result()
	if err != nil {
		log.Fatalf("unable to set task %v as processing: %v", task.Key, err)
	}

	log.Printf("SetTaskAsProcessing: set task %v to status: 'processing' with start time %v", task.Key, task.StartTime)
}

// SetTaskAsCompleteWithPassedResult will denote the current time and mark the task as TASKSTATUS_COMPLETE with TASKRESULT_PASSED
func (f RedisFacade) SetTaskAsCompleteWithPassedResult(task *model.TaskMetadata) {
	f.setTaskAsComplete(task, TASKRESULT_PASSED)
}

// SetTaskAsCompleteWithFailedResult will denote the current time and mark the task as TASKSTATUS_COMPLETE with TASKRESULT_FAILED
func (f RedisFacade) SetTaskAsCompleteWithFailedResult(task *model.TaskMetadata) {
	f.setTaskAsComplete(task, TASKRESULT_FAILED)
}

// SetTaskAsCompleteWithErrorResult will denote the current time and mark the task as TASKSTATUS_COMPLETE with TASKRESULT_ERROR
func (f RedisFacade) SetTaskAsCompleteWithErrorResult(task *model.TaskMetadata) {
	f.setTaskAsComplete(task, TASKRESULT_ERROR)
}

func (f RedisFacade) setTaskAsComplete(task *model.TaskMetadata, result string) {
	keys := GetTaskMetadataKeys()
	
	task.Status = TASKSTATUS_COMPLETE
	task.Result = result

	duration := time.Now().Sub(task.StartTime)
	task.Duration = int64(duration)
	task.DurationString = duration.Truncate(time.Second).String()

	values := []string{
		string(keys.Status),
		string(task.Status),
		string(keys.Result),
		string(task.Result),
		string(keys.Duration),
		fmt.Sprint(duration.Nanoseconds()),
		string(keys.DurationString),
		string(task.DurationString),
	}

	_, err := f.client.HSet(f.context, task.Key, values).Result()
	if err != nil {
		log.Fatalf("unable to set task %v as complete: %v", task.Key, err)
	}

	log.Printf("setTaskAsComplete: set task %v to status 'complete' with result '%v', duration %v", task.Key, task.Result, task.DurationString)
}