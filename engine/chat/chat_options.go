package chat

type MessageType string

const (
	AGETNT   MessageType = "agent"
	FEEDBACK MessageType = "feedback"
)

type Message struct {
	MessageType MessageType
	Response    *string
	Invocation  *Invocation
}

type ChatOption struct {
	systemPrompt string
	prompt       string
	history      []*Message
}

func NewChatOption(systemPrompt string, prompt string, history []*Message) *ChatOption {
	return &ChatOption{
		systemPrompt: systemPrompt,
		prompt:       prompt,
		history:      history,
	}
}

func (c *ChatOption) GetSystemPrompt() string {
	return c.systemPrompt
}

func (c *ChatOption) GetPrompt() string {
	return c.prompt
}

func (c *ChatOption) GetHistory() []*Message {
	return c.history
}

func (c *ChatOption) UpdateSystemPrompt(prompt string) {
	c.systemPrompt = prompt
}

func (c *ChatOption) UpdateHistroy(histroy []*Message) {
	c.history = histroy
}
