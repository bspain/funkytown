package main

import (
	"context"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/bspain/funkytown/shared/redisfacade"
)

const TASKQUEUE_CHECK_DELAY_IN_SECONDS = int64(5)
var worker_start_delay = time.Duration(TASKQUEUE_CHECK_DELAY_IN_SECONDS * time.Second.Nanoseconds())

var redis_host = os.Getenv("REDIS_HOST")
var redis_port = os.Getenv("REDIS_PORT")
var spec_root = os.Getenv("SPEC_ROOT")

var ctx = context.Background()

func main() {
	f := redisfacade.NewFacade(ctx, redis_host, redis_port)

	// wait until the run metadata is published - retry until connected
	var wait = true
	for wait {
		var errtype redisfacade.RedisFacadeErrorType
		_, err, errtype := f.GetRunMetadata()

		if err == nil {
			wait = false
		} else {
			switch errtype {
			case redisfacade.TcpHostErr:
				log.Printf("Waiting for redis controller (no such host), %v", err)
			case redisfacade.TcpCxErr:
				log.Printf("Waiting for redis controller (connection refused), %v", err)
			case redisfacade.RedisNil:
				log.Printf("Waiting for run metadata, %v", err)
			}
			time.Sleep(worker_start_delay)
		}
	}

	// Grab and item from the task queue
	var cont = true
	for cont {
		empty, key := f.PopTask()

		if empty {
			cont = false
			log.Printf("task queue was empty, worker exiting...")
			continue
		}

		// Get task and process
		t := f.GetTaskMetadata(key)

		f.SetTaskAsProcessing(&t)
		
		log.Printf("working on task: %v", t.Key)

		// TODO: Actual work
		work_time := rand.Intn(3)
		time.Sleep(time.Duration(int64(work_time) * time.Second.Nanoseconds()))

		f.SetTaskAsCompleteWithPassedResult(&t)
	}
}
