package serializer

import (
	"log"
	"testing"

	"github.com/runetale/notch/engine/chat"
	"github.com/runetale/notch/storage"
	"github.com/runetale/notch/types"
)

func Test_SerializeInvocation(t *testing.T) {
	payload := "Hello, World!"
	inv := chat.Invocation{
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
		log.Printf("Action: %s, Attributes: %v, Payload: %v\n",
			invocation.Action, invocation.Attributes, *invocation.Payload)
	}
}

func Test_ParseStorage(t *testing.T) {
	s := storage.NewStorage("Time Stroge", types.TIMER, nil)
	log.Println("Time Storage Output:")
	log.Println(paraseStorage(s))
}
