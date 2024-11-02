package llm

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type LLMTypeName string

const (
	Ollama    LLMTypeName = "ollama"
	OpenAI    LLMTypeName = "openai"
	Fireworks LLMTypeName = "fireworks"
	Groq      LLMTypeName = "groq"
)

type LLMOptions struct {
	typeName      LLMTypeName
	modelName     string
	contextWindow uint32
	host          string
	port          uint16
}

func NewLLMOptions(generator string, contextWindow uint32) (LLMOptions, error) {
	return parseGeneratorString(generator, contextWindow)
}

func parseGeneratorString(raw string, contextWindow uint32) (LLMOptions, error) {
	raw = strings.Trim(raw, " \"'")
	if raw == "" {
		return LLMOptions{}, errors.New("generator string can't be empty")
	}

	generator := LLMOptions{
		contextWindow: contextWindow,
	}

	localGeneratorPattern := `^([a-zA-Z0-9_]+)://([a-zA-Z0-9_-]+)@([a-zA-Z0-9_.-]+):(\d+)$`
	re := regexp.MustCompile(localGeneratorPattern)

	matches := re.FindStringSubmatch(raw)
	if matches == nil || len(matches) != 5 {
		return LLMOptions{}, fmt.Errorf("can't parse '%s' generator string", raw)
	}

	generator.typeName = LLMTypeName(matches[1])
	generator.modelName = matches[2]
	generator.host = matches[3]
	port, err := strconv.Atoi(matches[4])
	if err != nil {
		return LLMOptions{}, fmt.Errorf("invalid port: %s", matches[4])
	}
	generator.port = uint16(port)

	return generator, nil
}
