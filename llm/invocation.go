// invocation use by engine, parse `llm toolcalls` to invocation
package llm

import (
	"fmt"
	"strings"
)

type Invocation struct {
	Action     string
	Attributes map[string]string
	Payload    *string
}

func NewInvocation(
	action string,
	attributes map[string]string,
	payload *string,
) *Invocation {
	return &Invocation{
		Action:     action,
		Attributes: attributes,
		Payload:    payload,
	}
}

// showing by llm created commands to user prompt
func (i *Invocation) FunctionCallString() string {
	parts := []string{}
	if i.Payload != nil {
		parts = append(parts, *i.Payload)
	}

	if i.Attributes != nil {
		for name, value := range i.Attributes {
			parts = append(parts, fmt.Sprintf("%s=%s", name, value))
		}
	}
	return fmt.Sprintf("%s(%s)", i.Action, strings.Join(parts, ", "))
}
