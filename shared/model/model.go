package model

type RunMetadata struct {
	Key	string
	Runid string
	CommandCount	int
	CommandFinishedCount	int
	Finished bool
}

type RunMetadataKey string

const (
	KeyRunMetaId	RunMetadataKey = "runid"
	KeyRunMetaCommandCount RunMetadataKey = "cmdcount"
	KeyRunMetaCommandFinishedCount RunMetadataKey = "cmdfinishedcount"
	KeyRunMetaFinished RunMetadataKey = "finished"
)