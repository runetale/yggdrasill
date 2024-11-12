package events

import (
	"time"

	"github.com/runetale/notch/types"
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

type DisplayEvent interface {
	Display() string
}

type StateUpdateEvent struct {
	systemPrompt string
	prompt       string
	history      string
	savePath     string
}

func NewStateUpdateEvent(sys, prom, his, savePath string) DisplayEvent {
	return &StateUpdateEvent{
		systemPrompt: sys,
		prompt:       prom,
		history:      his,
		savePath:     savePath,
	}
}

func (e *StateUpdateEvent) Display() string {
	return ""
}

type InvalidResponseEvent struct {
	response string
}

func NewInvalidResponseEvent(res string) DisplayEvent {
	return &InvalidResponseEvent{
		response: res,
	}
}

func (e *InvalidResponseEvent) Display() string {
	return ""
}

type InvalidActionEvent struct {
	action string
	err    string
}

func NewInvalidActionEvent(ac, err string) DisplayEvent {
	return &InvalidActionEvent{
		action: ac,
		err:    err,
	}
}

func (e *InvalidActionEvent) Display() string {
	return ""
}

type ActionTimeoutEvent struct {
	action  string
	elapsed time.Duration
}

func NewActionTimeoutEvent(ac string, elapsed time.Duration) DisplayEvent {
	return &ActionTimeoutEvent{
		action:  ac,
		elapsed: elapsed,
	}
}

func (e *ActionTimeoutEvent) Display() string {
	return ""
}

type ActionExecutedEvent struct {
	invocation string
	err        *string
	result     *string
	elapsed    time.Duration
}

func NewActionExecutedEvent(inv string, err, result *string, elapsed time.Duration) DisplayEvent {
	return &ActionExecutedEvent{
		invocation: inv,
		err:        err,
		result:     result,
		elapsed:    elapsed,
	}
}

func (e *ActionExecutedEvent) Display() string {
	return ""
}

type TaskCompleteEvent struct {
	impossible bool
	reason     *string
}

func NewTaskCompleteEvent(impossible bool, reason *string) DisplayEvent {
	return &TaskCompleteEvent{
		impossible: impossible,
		reason:     reason,
	}
}

func (e *TaskCompleteEvent) Display() string {
	return ""
}

type StorageUpdateEvent struct {
	storageName string
	storageType types.StorageType
	key         string
	prev        *string
	new         *string
}

func NewStorageUpdateEvent(storageName string, t types.StorageType, key string, prev, new *string) DisplayEvent {
	return &StorageUpdateEvent{
		storageName: storageName,
		storageType: t,
		key:         key,
		prev:        prev,
		new:         new,
	}
}

func (e *StorageUpdateEvent) Display() string {
	return ""
}

type MetricsEvent struct {
	metrics string
}

func NewMetricsEvent(m string) DisplayEvent {
	return &MetricsEvent{
		metrics: m,
	}
}

func (e *MetricsEvent) Display() string {
	return ""
}

type EmptyResponseEvent struct {
}

func NewEmptyResponseEvent() DisplayEvent {
	return &EmptyResponseEvent{}
}

func (e *EmptyResponseEvent) Display() string {
	return "agent did not provide valid instructions: empty response"
}
