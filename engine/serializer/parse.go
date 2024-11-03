package serializer

import (
	"encoding/xml"
	"fmt"
	"log"
	"strings"

	"github.com/runetale/notch/llm"
)

type Parsed struct {
	Processed   int
	Invocations []*llm.Invocation
}

func preprocessBlock(ptr string) string {
	return ptr
}

func buildInvocation(closingName string, element xml.StartElement, payload string) (*llm.Invocation, error) {
	if element.Name.Local != closingName {
		return nil, fmt.Errorf("unexpected closing %s while parsing %s", closingName, element.Name.Local)
	}

	attributes := make(map[string]string)
	for _, attr := range element.Attr {
		attributes[attr.Name.Local] = attr.Value
	}

	return llm.NewInvocation(element.Name.Local, attributes, &payload), nil
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

		// XMLイベントの処理
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
