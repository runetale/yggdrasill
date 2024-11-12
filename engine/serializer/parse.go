// provided serializer support functions
package serializer

import (
	"encoding/xml"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/runetale/notch/engine/action"
	"github.com/runetale/notch/engine/chat"
	"github.com/runetale/notch/storage"
	"github.com/runetale/notch/types"
)

type Parsed struct {
	Processed   int
	Invocations []*chat.Invocation
}

func preprocessBlock(ptr string) string {
	return ptr
}

func buildInvocation(closingName string, element xml.StartElement, payload string) (*chat.Invocation, error) {
	if element.Name.Local != closingName {
		return nil, fmt.Errorf("unexpected closing %s while parsing %s", closingName, element.Name.Local)
	}

	attributes := make(map[string]string, 0)
	for _, attr := range element.Attr {
		attributes[attr.Name.Local] = attr.Value
	}

	return chat.NewInvocation(element.Name.Local, attributes, &payload), nil
}

func tryParseBlock(ptr string) Parsed {
	prevLen := len(ptr)
	ptr = preprocessBlock(ptr)
	delta := 0
	if len(ptr) != prevLen {
		delta = len(ptr) - prevLen
	}

	decoder := xml.NewDecoder(strings.NewReader(ptr))
	var parsed Parsed
	srcSize := len(ptr)

	var currElement *xml.StartElement
	var currPayload strings.Builder

	for {
		t, err := decoder.Token()
		if err != nil {
			break
		}

		switch event := t.(type) {
		case xml.StartElement:
			currElement = &event
			currPayload.Reset()
		case xml.CharData:
			currPayload.WriteString(string(event))
		case xml.EndElement:
			if currElement != nil {
				inv, err := buildInvocation(event.Name.Local, *currElement, currPayload.String())
				if err != nil {
					log.Println("Error:", err)
				} else {
					parsed.Invocations = append(parsed.Invocations, inv)
				}
				currElement = nil
			}
		default:
			log.Printf("Unexpected XML element: %v\n", event)
		}
	}

	srcSizeNow := len(ptr) - int(decoder.InputOffset())
	parsed.Processed = srcSize - srcSizeNow - delta

	return parsed
}

func parseInvocation(inv *chat.Invocation) string {
	var xml strings.Builder
	xml.WriteString(fmt.Sprintf("<%s", inv.Action))

	if inv.Attributes != nil {
		for key, value := range inv.Attributes {
			xml.WriteString(fmt.Sprintf(" %s=\"%s\"", key, value))
		}
	}

	payload := ""
	if inv.Payload != nil {
		payload = *inv.Payload
	}
	xml.WriteString(fmt.Sprintf(">%s</%s>", payload, inv.Action))

	return xml.String()
}

func parseAction(ac action.Action) string {
	var xml strings.Builder

	xml.WriteString(fmt.Sprintf("<%s", ac.Name()))

	attributes := ac.ExampleAttributes()
	for name, value := range attributes {
		xml.WriteString(fmt.Sprintf(` %s="%s"`, name, value))
	}

	if payload := ac.ExamplePayload(); payload != nil {
		xml.WriteString(fmt.Sprintf(">%s</%s>", *payload, ac.Name()))
	} else {
		xml.WriteString("/>")
	}

	return xml.String()
}

func paraseStorage(s *storage.Storage) string {
	if s.IsEmpty() {
		return ""
	}

	var result string

	switch s.GetStorageType() {
	case types.TIMER:
		startedAt := s.GetStartedAt()
		elapsed := time.Since(startedAt)
		result = fmt.Sprintf("## Current date: %s\n", time.Now().Format("01 January 2006 15:04"))
		result += fmt.Sprintf("## Time since start: %v\n", elapsed)

	case types.TAGGED:
		var xml strings.Builder
		xml.WriteString(fmt.Sprintf("<%s>\n", s.GetName()))
		for key, entry := range s.GetEntryList() {
			xml.WriteString(fmt.Sprintf("  - %s=%s\n", key, entry.Data))
		}
		xml.WriteString(fmt.Sprintf("</%s>", s.GetName()))
		result = xml.String()

	case types.UNTAGGED:
		var xml strings.Builder
		xml.WriteString(fmt.Sprintf("<%s>\n", s.GetName()))
		for _, entry := range s.GetEntries() {
			xml.WriteString(fmt.Sprintf("  - %s\n", entry.Data))
		}
		xml.WriteString(fmt.Sprintf("</%s>", s.GetName()))
		result = xml.String()

	case types.COMPLETION:
		var xml strings.Builder
		xml.WriteString(fmt.Sprintf("<%s>\n", s.GetName()))
		for _, entry := range s.GetEntries() {
			status := "not completed"
			if entry.Complete {
				status = "COMPLETED"
			}
			xml.WriteString(fmt.Sprintf("  - %s : %s\n", entry.Data, status))
		}
		xml.WriteString(fmt.Sprintf("</%s>", s.GetName()))
		result = xml.String()

	case types.CURRENTPREVIOUS:
		current, currentFound := s.GetEntry("CURRENT_TAG")
		if currentFound {
			result = fmt.Sprintf("* Current %s: %s", s.GetName(), strings.TrimSpace(current.Data))
			prev, prevFound := s.GetEntry("PREVIOUS_TAG")
			if prevFound {
				result += fmt.Sprintf("\n* Previous %s: %s", s.GetName(), strings.TrimSpace(prev.Data))
			}
		} else {
			result = ""
		}
	}

	return result
}
