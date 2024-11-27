// invocation use by engine, parsed `llm toolcalls` to invocation
// that's simply llm action
package chat

import (
	"fmt"
	"strings"

	"github.com/runetale/yggdrasill/engine/action"
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
	if len(attributes) == 0 {
		attributes = nil
	}
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

func (inv *Invocation) ValidateAction(ac action.Action) error {
	payloadRequired := ac.ExamplePayload() != nil
	attrsRequired := ac.ExampleAttributes() != nil
	hasPayload := inv.Payload != nil
	hasAttributes := inv.Attributes != nil

	if payloadRequired && !hasPayload {
		return fmt.Errorf("no content specified for '%s'", inv.Action)
	} else if attrsRequired && !hasAttributes {
		return fmt.Errorf("no attributes specified for '%s'", inv.Action)
	} else if !payloadRequired && hasPayload {
		return fmt.Errorf("no content needed for '%s'", inv.Action)
	} else if !attrsRequired && hasAttributes {
		return fmt.Errorf("no attributes needed for '%s'", inv.Action)
	}

	if attrsRequired {
		requiredAttrs := []string{}
		for key := range ac.ExampleAttributes() {
			requiredAttrs = append(requiredAttrs, key)
		}

		passedAttrs := []string{}
		for key := range inv.Attributes {
			passedAttrs = append(passedAttrs, key)
		}

		for _, required := range requiredAttrs {
			found := false
			for _, passed := range passedAttrs {
				if required == passed {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("no '%s' attribute specified for '%s'", required, inv.Action)
			}
		}
	}

	return nil
}
