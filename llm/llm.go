package llm

type LLMClinet struct {
}

type ChatOption struct {
	SystemPrompt string
	Prompt       string
	History      []string
}

func NewChatOption() *ChatOption {
	return &ChatOption{}
}
