// invocation use by engine, parse `llm toolcalls` to invocation
package llm

import (
	"fmt"
	"strings"
)

type Invocation struct {
	action     string
	attributes *map[string]string
	payload    *string
}

func NewInvocation(
	action string,
	attributes *map[string]string,
	payload *string,
) *Invocation {
	return &Invocation{
		action:     action,
		attributes: attributes,
		payload:    payload,
	}
}

// showing by llm created commands to user prompt
func (i *Invocation) FunctionCallString() string {
	parts := []string{}
	if i.payload != nil {
		parts = append(parts, *i.payload)
	}

	if i.attributes != nil {
		for name, value := range *i.attributes {
			parts = append(parts, fmt.Sprintf("%s=%s", name, value))
		}
	}
	return fmt.Sprintf("%s(%s)", i.action, strings.Join(parts, ", "))
}
