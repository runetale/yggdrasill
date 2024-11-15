package memory

import (
	_ "embed"
	"time"

	"github.com/runetale/yggdrasill/engine/action"
	"github.com/runetale/yggdrasill/storage"
	"github.com/runetale/yggdrasill/types"
)

//go:embed delete.prompt
var deletePrompt string

//go:embed ns.prompt
var nsPrompt string

type DeleteMemory struct {
}

func NewDeleteMemory() action.Action {
	return &DeleteMemory{}
}

func (m *DeleteMemory) Name() string {
	return "delete_memory"
}

func (m *DeleteMemory) Description() string {
	return deletePrompt
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
	return nsPrompt
}
