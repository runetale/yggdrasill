// openai, ollama, groq client
package llm

type LLMClientImpl interface {
	Chat(option *ChatOption) ([]*Invocation, string, error)
}

type LLMClinet struct {
}
