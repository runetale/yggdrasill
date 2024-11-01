package events

type Event string

const (
	MetricsUpdate Event = ""

	StorageUpdate Event = ""

	StateUpdate Event = ""

	InvalidUpdate Event = ""
	InvalidAction Event = ""

	ActionTimeOut  Event = ""
	ActionExecuted Event = ""

	TaskComplete Event = ""
)

// TODO Eventごとのチャネルと関数を設定する
type Channel struct {
	Sender   chan<- Event
	Receiver <-chan Event
}

func NewChannel() *Channel {
	ch := make(chan Event)
	return &Channel{
		Sender:   ch,
		Receiver: ch,
	}
}
