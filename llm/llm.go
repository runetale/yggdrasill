// openai, ollama, groq client
package llm

type LLMClientImpl interface {
	Chat(*ChatOption) (*Invocation, error)
}

type LLMClinet struct {
}
