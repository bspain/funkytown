package model

type RunMetadata struct {
	Key            string
	Runid          string
	TasksRemaining int
	TasksFinished  int
	Finished       bool
}

type TaskMetadata struct {
	KeyGroup string

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
