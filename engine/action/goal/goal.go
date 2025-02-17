package goal

import (
	_ "embed"
	"time"

	"github.com/runetale/notch/engine/action"
	"github.com/runetale/notch/storage"
	"github.com/runetale/notch/types"
)

//go:embed update.prompt
var updatePrompt string

//go:embed ns.prompt
var nsPrompt string

type Goal struct {
}

func NewGoal() action.Action {
	return &Goal{}
}

func (d *Goal) Name() string {
	return "update_goal"
}

func (d *Goal) Description() string {
	return updatePrompt
}

func (d *Goal) Run(storage *storage.Storage, attributes map[string]string, payload string) string {
	storage.SetCurrent(payload)
	return "goal updated"
}

func (d *Goal) Timeout() *time.Duration {
	return nil
}

func (d *Goal) ExamplePayload() *string {
	p := "your new goal"
	return &p
}

func (d *Goal) ExampleAttributes() map[string]string {
	return nil
}

func (d *Goal) RequiredVariables() []*string {
	return nil
}

func (d *Goal) RequiresUserConfirmation() bool {
	return true
}

func (d *Goal) GetNamespace() types.NamespaceType {
	return types.GOAL
}

func (d *Goal) NamespaceDescription() string {
	return nsPrompt
}
