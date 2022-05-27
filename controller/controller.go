package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bspain/funkytown/controller/lib/groupedtaskloader"
	"github.com/bspain/funkytown/shared/model"
	"github.com/bspain/funkytown/shared/redisfacade"
)

const TASK_SCANNING_LOOP_SLEEP_TIME_IN_SECONDS = 10

var scan_loop_sleep = time.Duration(TASK_SCANNING_LOOP_SLEEP_TIME_IN_SECONDS * time.Second.Nanoseconds())
var redis_host = os.Getenv("REDIS_HOST")
var redis_port = os.Getenv("REDIS_PORT")
var grouptasksfile = os.Getenv("GROUP_TASKS_FILE")
var html_index_file = os.Getenv("HTML_INDEX_FILE")

var ctx = context.Background()

func main() {
	log.Printf("funkytown Controller has started...")
	f := redisfacade.NewFacade(ctx, redis_host, redis_port)

	tasks := ParseTaskFile(grouptasksfile)
	PushTasks(f, tasks)
	go MonitorQueue(f, tasks)
	HostReporter(f, tasks)
}

// ParseTaskFile parses the group tasks file and returns a map of all tasks
func ParseTaskFile(file string) map[string]model.Task {
	// Load in tasks file
	groupedtasklist, err := groupedtaskloader.LoadGroupTasksFile(file)
	if err != nil {
		panic(err)
	}

	size := groupedtasklist.TaskCount()
	tasks := make(map[string]model.Task, size)

	for _, group := range groupedtasklist.Groups {
		for i, task := range group.Tasks {

			key := fmt.Sprintf("task:%v:%v", group.Name, i)
			tasks[key] = task
		}
	}

	return tasks
}

func PushTasks(f redisfacade.RedisFacade, tasks map[string]model.Task) {
	for key, task := range tasks {
		// Create a task metadata object for the task
		f.SetTaskMetadata(key, task)

		// Push the task to the workqueue
		f.PushTask(key)
	}

	size := len(tasks)
	// Set the run metadata, this is the signal to the workers that work is ready to begin.
	f.SetRunMetadata("a_new_run", size)
}

func MonitorQueue(f redisfacade.RedisFacade, tasks map[string]model.Task) {

	// Remap all tasks to key/status
	size := len(tasks)
	finished_tasks := make(map[string]bool, size)
	for key := range tasks {
		finished_tasks[key] = false
	}

	// Monitor tasks loop
	var t_complete = 0
	for t_complete < size {
		var t_waiting, t_processing = 0, 0
		
		for key, done := range finished_tasks {
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
				finished_tasks[key] = true
				t_complete++
			}
		}

		log.Printf("task scan complete: %v waiting, %v processing, %v complete", t_waiting, t_processing, t_complete)
		time.Sleep(scan_loop_sleep)
	}

	// TODO: Set run metadata finished
	log.Printf("all tasks complete.  Monitor exiting...")
}

type CombinedMetaAndResults struct {
	Meta model.RunMetadata
	Tasks []model.TaskMetadata
}

func HostReporter(f redisfacade.RedisFacade, tasks map[string]model.Task) {
	log.Printf("Starting reporting server...")

	template := template.Must(template.ParseFiles(html_index_file))
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/results", func(w http.ResponseWriter, r *http.Request) {
		meta, err, _ := f.GetRunMetadata()

		var t_metas []model.TaskMetadata
		for task := range tasks {
			t_meta := f.GetTaskMetadata(task)
			t_metas = append(t_metas, t_meta)
		}

		cmr := CombinedMetaAndResults{
			Meta: meta,
			Tasks: t_metas,
		}
		if err != nil {
			log.Printf("Error fetching run metadata: %v", err)
			http.Error(w, "Something went wrong!", 500)
			return
		}
		template.Execute(w, cmr)
	})

	http.ListenAndServe(":3000", nil)
}