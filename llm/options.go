package llm

type ChatOption struct {
	SystemPrompt string
	Prompt       string
	History      []string
}

func NewChatOption(systemPrompt string, prompt string, history []string) *ChatOption {
	return &ChatOption{
		SystemPrompt: systemPrompt,
		Prompt:       prompt,
		History:      history,
	}
}
