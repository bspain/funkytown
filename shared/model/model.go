package model

import "time"

type RunMetadata struct {
	Key            string
	Runid          string
	TasksRemaining int
	TasksFinished  int
	Finished       bool
}

type TaskMetadata struct {
	Key string
	Spec string
	Viewport string
	Browser string
	Status string
	Result string
	StartTime time.Time
	Duration int64
	DurationString string
}

type Task struct {
	Spec	string
	Viewport	string
	Browser	string
}

type Group struct {
	Name	string
	Tasks	[]Task
}

type GroupedTasks struct {
	Groups []Group
}

func (gt GroupedTasks) TaskCount() int {
	var count = 0
	for _, g := range gt.Groups {
		count = count + len(g.Tasks)
	}

	return count
}
