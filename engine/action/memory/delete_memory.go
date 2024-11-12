package memory

import (
	"os"
	"time"

	"github.com/runetale/notch/engine/action"
	"github.com/runetale/notch/storage"
	"github.com/runetale/notch/types"
)

type DeleteMemory struct {
}

func NewDeleteMemory() action.Action {
	return &DeleteMemory{}
}

func (m *DeleteMemory) Name() string {
	return "delete_memory"
}

func (m *DeleteMemory) Description() string {
	filepath := "delete.prompt"
	data, _ := os.ReadFile(filepath)
	return string(data)
}

func (m *DeleteMemory) Run(storage *storage.Storage, attributes map[string]string, payload string) string {
	key := attributes["key"]
	storage.AddTagged(key, payload)
	return "memory saved"
}

func (m *DeleteMemory) Timeout() *time.Duration {
	return nil
}

func (m *DeleteMemory) ExamplePayload() *string {
	p := ""
	return &p
}

func (m *DeleteMemory) ExampleAttributes() map[string]string {
	attr := map[string]string{}
	attr["key"] = "note"
	return attr
}

func (m *DeleteMemory) RequiredVariables() []*string {
	return nil
}

func (m *DeleteMemory) RequiresUserConfirmation() bool {
	return true
}

func (m *DeleteMemory) GetNamespace() types.NamespaceType {
	return types.MEMORY
}

func (m *DeleteMemory) NamespaceDescription() string {
	filepath := "ns.prompt"
	data, _ := os.ReadFile(filepath)
	return string(data)
}
