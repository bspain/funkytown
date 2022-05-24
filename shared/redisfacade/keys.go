package redisfacade

type RunMetadataKeys struct {
	RunID 	string
	TasksRemaining string
	TasksFinished	string
	Finished	string
}

func GetRunMetadataKeys() RunMetadataKeys {
	return RunMetadataKeys{
		"runid",
		"tasksremaining",
		"tasksfinished",
		"finished",
	}
}

type TaskMetadataKeys struct {
	Group string
	Spec string
	Viewport string
	Browser string
	Status string
	Result string
	StartTime string
	Iterations string
	Duration string
	DurationString string
}

func GetTaskMetadataKeys() TaskMetadataKeys {
	return TaskMetadataKeys{
		"group",
		"spec",
		"viewport",
		"browser",
		"status",
		"result",
		"starttime",
		"iterations",
		"duration",
		"durationstring",
	}
}
