// xml pareser for state
package serializer

import (
	"bytes"
	"encoding/xml"
	"html/template"

	"github.com/runetale/notch/engine/state"
)

type SystemPrompt struct {
	SystemPrompt     string `xml:"system_prompt"`
	Storages         string `xml:"storages"`
	Iterations       string `xml:"iterations"`
	AvailableActions string `xml:"available_actions"`
	Guidance         string `xml:"guidance"`
}

func displaySystemPrompt(state *state.State) (string, error) {
	// input data to template
	tmpl, err := template.New("prompt").ParseFiles("system.prompt")
	if err != nil {
		return "", err
	}

	var result bytes.Buffer
	err = tmpl.Execute(&result, map[string]string{
		"SystemPrompt":     systemPrompt,
		"Storages":         string(serializedStorages),
		"Iterations":       iterations,
		"AvailableActions": availableActions,
		"Guidance":         guidance,
	})
	if err != nil {
		return "", err
	}

	// parsing xml
	promptXML := SystemPrompt{
		SystemPrompt:     systemPrompt,
		Storages:         string(serializedStorages),
		Iterations:       iterations,
		AvailableActions: availableActions,
		Guidance:         guidance,
	}

	xmlOutput, err := xml.MarshalIndent(promptXML, "", "  ")
	if err != nil {
		return "", err
	}

	return string(xmlOutput), nil
}
