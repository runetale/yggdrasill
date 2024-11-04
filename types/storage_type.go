package types

type StorageType string

const (
	UNTAGGED        = "untagged"
	TAGGED          = "tagged"
	CURRENTPREVIOUS = "current-previous"
	COMPLETION      = "completion"
	TIMER           = "timer"
)
