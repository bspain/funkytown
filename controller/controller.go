package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/bspain/funkytown/controller/lib/groupedtaskloader"
	"github.com/bspain/funkytown/shared/redisfacade"
)

const TASK_SCANNING_LOOP_SLEEP_TIME_IN_SECONDS = 10

var scan_loop_sleep = time.Duration(TASK_SCANNING_LOOP_SLEEP_TIME_IN_SECONDS * time.Second.Nanoseconds())
var redis_host = os.Getenv("REDIS_HOST")
var redis_port = os.Getenv("REDIS_PORT")
var grouptasksfile = os.Getenv("GROUP_TASKS_FILE")

var ctx = context.Background()

func main() {
	log.Printf("funkytown Controller has started...")

	// Load in tasks file
	groupedtasklist, err := groupedtaskloader.LoadGroupTasksFile(grouptasksfile)
	if err != nil {
		panic(err)
	}

	f := redisfacade.NewFacade(ctx, redis_host, redis_port)
	size := groupedtasklist.TaskCount()

	f.SetRunMetadata("a_new_run", size)

	tasks := make(map[string]bool, size)

	for _, group := range groupedtasklist.Groups {
		for i, task := range group.Tasks {
			// Create a task metadata object for the task
			key := f.SetTaskMetadata(group.Name, i, task)

			// Push the task to the workqueue
			f.PushTask(key)

			// Record the key for the monitoring loop
			tasks[key] = false
		}
	}

	// Monitor tasks loop
	var t_complete = 0
	for t_complete < size {
		var t_waiting, t_processing = 0, 0
		
		for key, done := range tasks {
			// Skip if we know the task is already finished
			if done == true {
				continue
			}

			t_status := f.GetTaskStatus(key)
			if t_status == redisfacade.TASKSTATUS_WAITING {
				t_waiting++
				continue
			}

			if t_status == redisfacade.TASKSTATUS_PROCESSING {
				t_processing++
				continue
				// TODO: Check for timeouts
			}
			if t_status == redisfacade.TASKSTATUS_COMPLETE {
				tasks[key] = true
				t_complete++
			}
		}

		log.Printf("task scan complete: %v waiting, %v processing, %v complete", t_waiting, t_processing, t_complete)
		time.Sleep(scan_loop_sleep)
	}

	log.Printf("all tasks complete.  Controller exiting...")
}