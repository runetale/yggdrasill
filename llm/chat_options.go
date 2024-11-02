package llm

type MessageType string

const (
	AGETNT   MessageType = "agent"
	FEEDBACK MessageType = "feedback"
)

type Message struct {
	messageType MessageType
	response    string
	invocation  Invocation
}

type ToolCall struct {
	id       string
	function Function
	theType  string
}

type Function struct {
	name string
	args string
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
