package llm

import (
	"fmt"
	"testing"
)

func Test_Parse(t *testing.T) {
	raw := "openai://gpt-4@localhost:12321"
	contextWindow := 2048

	generator, err := parseGeneratorString(raw, uint32(contextWindow))
	if err != nil {
		t.Fatalf(err.Error())
	}

	fmt.Printf("TypeName: %s\n", generator.typeName)
	fmt.Printf("ModelName: %s\n", generator.modelName)
	fmt.Printf("Host: %s\n", generator.host)
	fmt.Printf("Port: %d\n", generator.port)
	fmt.Printf("ContextWindow: %d\n", generator.contextWindow)
}
