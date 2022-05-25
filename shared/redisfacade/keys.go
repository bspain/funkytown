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
	Key string
	Group string
	Spec string
	Viewport string
	Browser string
	Status string
	Result string
	StartTime string
	Duration string
	DurationString string
}

func GetTaskMetadataKeys() TaskMetadataKeys {
	return TaskMetadataKeys{
		"key",
		"group",
		"spec",
		"viewport",
		"browser",
		"status",
		"result",
		"starttime",
		"duration",
		"durationstring",
	}
}
