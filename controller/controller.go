package main

import (
	"context"
	"log"
	"os"

	"github.com/bspain/funkytown/shared/redisfacade"
	"github.com/bspain/funkytown/controller/lib/groupedtaskloader"
)

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

	// Load in all tasks into the queue
	f.SetRunMetadata("a_new_run", groupedtasklist.TaskCount())

	for _, group := range groupedtasklist.Groups {
		for i, task := range group.Tasks {
			f.SetTaskMetadata(group.Name, i, task)
		}
	}	
}