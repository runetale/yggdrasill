package goal

import (
	"os"
	"time"

	"github.com/runetale/notch/engine/action"
	"github.com/runetale/notch/storage"
	"github.com/runetale/notch/types"
)

type Goal struct {
	storageType types.StorageType
	predefined  *map[string]string
}

func NewGoal() action.Action {
	return &Goal{
		storageType: types.UNTAGGED,
		predefined:  nil,
	}
}

func (d *Goal) StorageType() types.StorageType {
	return d.storageType
}

func (d *Goal) Predefined() *map[string]string {
	return d.predefined
}

func (d *Goal) Name() string {
	return "update_goal"
}

func (d *Goal) Description() string {
	filepath := "update.prompt"
	data, _ := os.ReadFile(filepath)
	return string(data)
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
