package types

type StorageType string

const (
	UNTAGGED        StorageType = "untagged"
	TAGGED          StorageType = "tagged"
	CURRENTPREVIOUS StorageType = "current-previous"
	COMPLETION      StorageType = "completion"
	TIMER           StorageType = "timer"
)
