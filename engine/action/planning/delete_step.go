package planning

import (
	"os"
	"strconv"
	"time"

	"github.com/runetale/notch/engine/action"
	"github.com/runetale/notch/storage"
	"github.com/runetale/notch/types"
)

type DeleteStep struct {
	storageType types.StorageType
	predefined  *map[string]string
}

func NewDeleteStep() action.Action {
	return &DeleteStep{
		storageType: types.UNTAGGED,
		predefined:  nil,
	}
}

func (a *DeleteStep) StorageType() types.StorageType {
	return a.storageType
}

func (a *DeleteStep) Predefined() *map[string]string {
	return a.predefined
}

func (a *DeleteStep) Name() string {
	return "delete_plan_step"
}

func (a *DeleteStep) Description() string {
	filepath := "delete.prompt"
	data, _ := os.ReadFile(filepath)
	return string(data)
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
	return types.GOAL
}

func (a *DeleteStep) NamespaceDescription() string {
	filepath := "ns.prompt"
	data, _ := os.ReadFile(filepath)
	return string(data)
}
