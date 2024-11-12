package events

type Event struct {
	eventType                 EventType
	name                      string
	happened                  string
	stateUpdateEvent          StateUpdateEvent
	invalidResonseEvent       InvalidResonseEvent
	invalidActionTimeoutEvent InvalidActionTimeoutEvent
	actionExecutedEvent       ActionExecutedEvent
	taskCompleteEvent         TaskCompleteEvent
	storageUpdateEvent        StorageUpdateEvent
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

func (e *Event) Name() string {
	return e.name
}

func (e *Event) Happened() string {
	return e.happened
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
