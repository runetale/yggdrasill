package llm

import (
	"log"
	"testing"
)

func Test_Parse(t *testing.T) {
	raw := "openai://gpt-4@localhost:12321"
	contextWindow := 2048

	generator, err := parseGeneratorString(raw, uint32(contextWindow))
	if err != nil {
		t.Fatalf(err.Error())
	}

	log.Printf("TypeName: %s\n", generator.typeName)
	log.Printf("ModelName: %s\n", generator.modelName)
	log.Printf("Host: %s\n", generator.host)
	log.Printf("Port: %d\n", generator.port)
	log.Printf("ContextWindow: %d\n", generator.contextWindow)
}
