package planning

import (
	_ "embed"
	"time"

	"github.com/runetale/yggdrasill/engine/action"
	"github.com/runetale/yggdrasill/storage"
	"github.com/runetale/yggdrasill/types"
)

//go:embed add.prompt
var addPrompt string

//go:embed ns.prompt
var nsPrompt string

type AddStep struct {
}

func NewAddStep() action.Action {
	return &AddStep{}
}

func (a *AddStep) Name() string {
	return "add_plan_step"
}

func (a *AddStep) Description() string {
	return addPrompt
}

func (a *AddStep) Run(storage *storage.Storage, attributes map[string]string, payload string) string {
	storage.AddCompletion(payload)
	return "step added to the plan"
}

func (a *AddStep) Timeout() *time.Duration {
	return nil
}

func (a *AddStep) ExamplePayload() *string {
	p := "complete the task"
	return &p
}

func (a *AddStep) ExampleAttributes() map[string]string {
	return nil
}

func (a *AddStep) RequiredVariables() []*string {
	return nil
}

func (a *AddStep) RequiresUserConfirmation() bool {
	return true
}

func (a *AddStep) GetNamespace() types.NamespaceType {
	return types.PLANNING
}

func (a *AddStep) NamespaceDescription() string {
	return nsPrompt
}
