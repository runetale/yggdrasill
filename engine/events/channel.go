package events

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

type Event struct {
	eventType EventType
	name      string
	happened  string
}

func NewEvent(evenType EventType, name, happened string) *Event {
	return &Event{
		eventType: evenType,
		name:      name,
		happened:  happened,
	}
}

func (e *Event) EventType() EventType {
	return e.eventType
}

type Channel struct {
	Chan chan *Event
}

func NewChannel() *Channel {
	ch := make(chan *Event)
	return &Channel{
		Chan: ch,
	}
}
