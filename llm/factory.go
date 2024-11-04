// openai, ollama, groq client
package llm

import "errors"

type LLMClientImpl interface {
	Chat(option *ChatOption) ([]*Invocation, string)
}

type LLMFactory struct {
	typeName      LLMTypeName
	modelName     string
	contextWindow uint32
	host          string
	port          uint16

	client LLMClientImpl
}

func NewLLMFactory(options LLMOptions, apiKey string) (*LLMFactory, error) {
	client, err := newLLMFactory(options.typeName, options, apiKey)
	if err != nil {
		return nil, err
	}

	return &LLMFactory{
		typeName:      options.typeName,
		modelName:     options.modelName,
		contextWindow: options.contextWindow,
		host:          options.host,
		port:          options.port,
		client:        client,
	}, nil
}

func newLLMFactory(llmType LLMTypeName, options LLMOptions, apiKey string) (LLMClientImpl, error) {
	switch llmType {
	case Ollama:
		return nil, errors.New("not implement ollama")
	case OpenAI:
		return NewOpenAIClient(options.modelName, apiKey, options.host, options.port), nil
	case Fireworks:
		return nil, errors.New("not implement fireworks")
	case Groq:
		return nil, errors.New("not implement groq")
	}
	return nil, errors.New("not suuported llm")
}

func (c *LLMFactory) Chat(options *ChatOption) ([]*Invocation, string) {
	return c.client.Chat(options)
}
