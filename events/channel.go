package events

type Channel struct {
	Chan chan DisplayEvent
}

func NewChannel() *Channel {
	ch := make(chan DisplayEvent)
	return &Channel{
		Chan: ch,
	}
}
