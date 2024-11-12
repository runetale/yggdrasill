package events

import (
	"fmt"
	"os"
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

func NewStateUpdateEvent(sys, prom, his string, savePath string) DisplayEvent {
	return &StateUpdateEvent{
		systemPrompt: sys,
		prompt:       prom,
		history:      his,
		savePath:     savePath,
	}
}

func (e *StateUpdateEvent) Display() string {
	data := ""
	if e.savePath != "" {
		data = fmt.Sprintf(
			"[SYSTEM PROMPT]\n\n%s\n\n[PROMPT]\n\n%s\n\n[CHAT]\n\n%s",
			e.systemPrompt,
			e.prompt,
			e.history,
		)

		err := os.WriteFile(e.savePath, []byte(data), 0644)
		if err != nil {
			return fmt.Sprintf("Error writing to %s: %v", e.savePath, err)
		}
	}
	return "updated state"
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
	return fmt.Sprintf("agent did not provide valid instructions\n\n%s\n\n", e.response)
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
	return fmt.Sprintf("invalid action %s : %s", e.action, e.err)
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
	return fmt.Sprintf("action %s timed out after %s", e.action, e.elapsed.String())
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
	if e.err != nil {
		return fmt.Sprintf("%s: %s", e.invocation, *e.err)
	}
	if e.result != nil {
		return fmt.Sprintf("%s -> %s bytes in %s", e.invocation, *e.result, e.elapsed.String())
	}
	return fmt.Sprintf("%s %s in %s", e.invocation, "no output", e.elapsed.String())
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
	reason := ""
	if e.impossible {
		if e.reason != nil {
			reason = *e.reason
		} else {
			reason = "no reason provided"
		}
		return fmt.Sprintf("task is impossible: %s", reason)
	}

	if e.reason != nil {
		reason = *e.reason
	} else {
		reason = "no reason provided"
	}
	return fmt.Sprintf("task complete: %s", reason)
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
	if e.prev == nil && e.new == nil {
		return fmt.Sprintf("storage.%s cleared", e.storageName)
	}
	if e.prev != nil && e.new == nil {
		return fmt.Sprintf("%s.%s removed", e.storageName, e.key)
	}
	if e.new != nil {
		return fmt.Sprintf("%s.%s > %s", e.storageName, e.key, *e.new)
	}
	return fmt.Sprintf("%s.%s prev=%s new=%s", e.storageName, e.key, *e.prev, *e.new)
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
	return e.metrics
}

type EmptyResponseEvent struct {
}

func NewEmptyResponseEvent() DisplayEvent {
	return &EmptyResponseEvent{}
}

func (e *EmptyResponseEvent) Display() string {
	return "agent did not provide valid instructions: empty response"
}
