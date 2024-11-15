package planning

import (
	_ "embed"
	"strconv"
	"time"

	"github.com/runetale/notch/engine/action"
	"github.com/runetale/notch/storage"
	"github.com/runetale/notch/types"
)

//go:embed delete.prompt
var deletePrompt string

type DeleteStep struct {
}

func NewDeleteStep() action.Action {
	return &DeleteStep{}
}

func (a *DeleteStep) Name() string {
	return "delete_plan_step"
}

func (a *DeleteStep) Description() string {
	return deletePrompt
}

func (a *DeleteStep) Run(storage *storage.Storage, attributes map[string]string, payload string) string {
	pos, _ := strconv.Atoi(payload)
	storage.DelCompletion(pos)
	return "step removed from the plan"
}

func (a *DeleteStep) Timeout() *time.Duration {
	return nil
}

func (a *DeleteStep) ExamplePayload() *string {
	p := "2"
	return &p
}

func (a *DeleteStep) ExampleAttributes() map[string]string {
	return nil
}

func (a *DeleteStep) RequiredVariables() []*string {
	return nil
}

func (a *DeleteStep) RequiresUserConfirmation() bool {
	return true
}

func (a *DeleteStep) GetNamespace() types.NamespaceType {
	return types.PLANNING
}

func (a *DeleteStep) NamespaceDescription() string {
	return nsPrompt
}
