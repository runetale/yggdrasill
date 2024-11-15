package planning

import (
	_ "embed"
	"time"

	"github.com/runetale/notch/engine/action"
	"github.com/runetale/notch/storage"
	"github.com/runetale/notch/types"
)

//go:embed clear.prompt
var clearPrompt string

type Clear struct {
}

func NewClear() action.Action {
	return &Clear{}
}

func (a *Clear) Name() string {
	return "clear_plan"
}

func (a *Clear) Description() string {
	return clearPrompt
}

func (a *Clear) Run(storage *storage.Storage, attributes map[string]string, payload string) string {
	return "plan clear"
}

func (a *Clear) Timeout() *time.Duration {
	return nil
}

func (a *Clear) ExamplePayload() *string {
	p := "complete the task"
	return &p
}

func (a *Clear) ExampleAttributes() map[string]string {
	return nil
}

func (a *Clear) RequiredVariables() []*string {
	return nil
}

func (a *Clear) RequiresUserConfirmation() bool {
	return true
}

func (a *Clear) GetNamespace() types.NamespaceType {
	return types.PLANNING
}

func (a *Clear) NamespaceDescription() string {
	return nsPrompt
}
