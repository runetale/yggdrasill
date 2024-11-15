package memory

import (
	_ "embed"
	"time"

	"github.com/runetale/notch/engine/action"
	"github.com/runetale/notch/storage"
	"github.com/runetale/notch/types"
)

//go:embed save.prompt
var savePrompt string

type SaveMemory struct {
}

func NewSaveMemroy() action.Action {
	return &SaveMemory{}
}

func (m *SaveMemory) Name() string {
	return "save_memory"
}

func (m *SaveMemory) Description() string {
	return savePrompt
}

func (m *SaveMemory) Run(storage *storage.Storage, attributes map[string]string, payload string) string {
	key := attributes["key"]
	storage.AddTagged(key, payload)
	return "memory saved"
}

func (m *SaveMemory) Timeout() *time.Duration {
	return nil
}

func (m *SaveMemory) ExamplePayload() *string {
	p := "your new goal"
	return &p
}

func (m *SaveMemory) ExampleAttributes() map[string]string {
	attr := map[string]string{}
	attr["key"] = "note"
	return attr
}

func (m *SaveMemory) RequiredVariables() []*string {
	return nil
}

func (m *SaveMemory) GetNamespace() types.NamespaceType {
	return types.MEMORY
}

func (m *SaveMemory) RequiresUserConfirmation() bool {
	return true
}

func (m *SaveMemory) NamespaceDescription() string {
	return nsPrompt
}
