// xml pareser for state
package serializer

import (
	"bytes"
	_ "embed"
	"fmt"
	"html"
	"html/template"
	"sort"
	"strings"

	"github.com/runetale/yggdrasill/engine/action"
	"github.com/runetale/yggdrasill/engine/chat"
	"github.com/runetale/yggdrasill/engine/state"
	"github.com/runetale/yggdrasill/storage"
)

//go:embed actions.prompt
var actionPrompt string

//go:embed system.prompt
var systemPrompt string

type System struct {
	SystemPrompt     string
	Storages         string
	Iterations       string
	AvailableActions string
	Guidance         string
}

func DisplaySystemPrompt(state *state.State) (string, error) {
	// input data to template
	tmpl, err := template.New("prompt").Parse(systemPrompt)
	if err != nil {
		return "", err
	}

	// system prompt
	task := state.GetTask()
	sysprompt := task.GetSystemPrompt()

	// storages
	// sort
	storages := state.GetStorages()
	sortedStorageKeys := make([]string, 0, len(storages))
	for key := range storages {
		sortedStorageKeys = append(sortedStorageKeys, key)
	}
	sort.Strings(sortedStorageKeys)

	// serialization
	serializedStorage := []string{}
	for _, key := range sortedStorageKeys {
		serializedStorage = append(serializedStorage, serializeStorage(storages[key]))
	}
	displayStorages := strings.Join(serializedStorage, "\n\n")

	// guidance
	guidance := strings.Join(task.GetGuidance(), "\n")

	// available actions
	actions, err := actionsForState(state)
	if err != nil {
		return "", err
	}
	availableActions := actionPrompt + "\n" + actions

	// iterations
	iterations := ""
	if state.GetMaxIteration() > 0 {
		iterations = fmt.Sprintf("You are currently at step %d of a maximum of %d", state.GetCurrentStep()+1, state.GetMaxIteration())
	}

	data := System{
		SystemPrompt:     sysprompt,
		Storages:         displayStorages,
		Iterations:       iterations,
		AvailableActions: availableActions,
		Guidance:         guidance,
	}

	var output bytes.Buffer
	if err := tmpl.Execute(&output, data); err != nil {
		return "", err
	}

	return output.String(), nil
}

func serializeAction(ac action.Action) string {
	var builder strings.Builder

	// create xml tag
	builder.WriteString(fmt.Sprintf("<%s", html.EscapeString(ac.Name())))

	// if existing attributes by aciton
	for name, exampleValue := range ac.ExampleAttributes() {
		escapedName := html.EscapeString(name)
		escapedValue := html.EscapeString(exampleValue)
		builder.WriteString(fmt.Sprintf(` %s="%s"`, escapedName, escapedValue))
	}

	// if existing payload by aciton
	if payload := ac.ExamplePayload(); payload != nil {
		builder.WriteString(fmt.Sprintf(
			">%s</%s>",
			html.EscapeString(*payload),
			html.EscapeString(ac.Name()),
		))
	} else {
		builder.WriteString("/>")
	}

	return builder.String()
}

func actionsForState(state *state.State) (string, error) {
	var builder strings.Builder

	for _, group := range state.GetNamespaces() {
		builder.WriteString(fmt.Sprintf("## %s\n\n", group.Name()))
		if group.Description() != "" {
			builder.WriteString(fmt.Sprintf("%s\n\n", group.Description()))
		}

		for _, action := range group.GetActions() {
			builder.WriteString(fmt.Sprintf("%s %s\n\n",
				action.Description(),
				serializeAction(action),
			))
		}
	}

	return builder.String(), nil
}

func SerializeInvocation(inv *chat.Invocation) *string {
	invocation := parseInvocation(inv)
	return &invocation
}

func SerializeAction(ac action.Action) string {
	return parseAction(ac)
}

func serializeStorage(s *storage.Storage) string {
	return paraseStorage(s)
}

func TryParse(raw string) []*chat.Invocation {
	ptr := raw
	var parsedInvocations []*chat.Invocation
	uniqueMap := make(map[string]bool)

	for {
		openIdx := strings.Index(ptr, "<")
		if openIdx == -1 {
			break
		}

		ptr = ptr[openIdx:]

		parsedBlock := tryParseBlock(ptr)
		if parsedBlock.Processed == 0 {
			break
		}

		for _, inv := range parsedBlock.Invocations {
			uniqueKey := fmt.Sprintf("%s-%v-%v", inv.Action, inv.Attributes, inv.Payload)
			if !uniqueMap[uniqueKey] {
				uniqueMap[uniqueKey] = true
				parsedInvocations = append(parsedInvocations, inv)
			}
		}

		ptr = ptr[parsedBlock.Processed:]
	}

	return parsedInvocations
}
