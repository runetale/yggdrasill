package serializer

import (
	"fmt"
	"testing"

	"github.com/runetale/notch/llm"
)

func Test_SerializeInvocation(t *testing.T) {
	payload := "Hello, World!"
	inv := llm.Invocation{
		Action: "greeting",
		Attributes: map[string]string{
			"lang": "en",
			"tone": "friendly",
		},
		Payload: &payload,
	}

	result := SerializeInvocation(&inv)
	t.Log(*result)
}

func Test_TryParse(t *testing.T) {
	raw := `<example><tag>Hello</tag><tag>World</tag></example>`
	invocations := TryParse(raw)

	for _, invocation := range invocations {
		fmt.Printf("Action: %s, Attributes: %v, Payload: %v\n",
			invocation.Action, invocation.Attributes, *invocation.Payload)
	}
}
