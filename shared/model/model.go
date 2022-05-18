package model

type RunMetadata struct {
	Key            string
	Runid          string
	TasksRemaining int
	TasksFinished  int
	Finished       bool
}

type RunMetadataKey string

const (
	KeyRunMetaId             RunMetadataKey = "runid"
	KeyRunMetaTasksRemaining RunMetadataKey = "tasksremaining"
	KeyRunMetaTasksFinished  RunMetadataKey = "tasksfinished"
	KeyRunMetaFinished       RunMetadataKey = "finished"
)
