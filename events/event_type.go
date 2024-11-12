package events

import (
	"time"
)

type EventType string

const (
	MetricsUpdate   EventType = "metrics_update"
	StorageUpdate   EventType = "storage_update"
	StateUpdate     EventType = "state_update"
	InvalidUpdate   EventType = "invaild_update"
	InvalidAction   EventType = "invalid_action"
	InvalidResponse EventType = "invalid_response"
	ActionTimeOut   EventType = "action_timeout"
	ActionExecuted  EventType = "aciton_executed"
	TaskComplete    EventType = "task_comlete"
	EmptyResponse   EventType = "empty_response"
)

type StateUpdateEvent struct {
	SystemPrompt string
	Prompt       string
	History      string
}

type InvalidResonseEvent struct {
	Response string
}

type InvalidActionEvent struct {
	Action string
	Error  string
}

type InvalidActionTimeoutEvent struct {
	Action  string
	Elapsed time.Duration
}

type ActionExecutedEvent struct {
	Invocation string
	Error      *string
	Result     *string
	Elapsed    time.Duration
}

type TaskCompleteEvent struct {
	Impossible bool
	Reason     *string
}

type StorageUpdateEvent struct {
	StorageName string
	Key         string
	Prev        *string
	New         *string
}
