package planning

import (
	_ "embed"
	"fmt"
	"strconv"
	"time"

	"github.com/runetale/notch/engine/action"
	"github.com/runetale/notch/storage"
	"github.com/runetale/notch/types"
)

//go:embed set-incomplete.prompt
var setInCompletePrompt string

type SetInComplete struct {
}

func NewSetInComplete() action.Action {
	return &SetInComplete{}
}

func (s *SetInComplete) Name() string {
	return "set_step_incompleted"
}

func (s *SetInComplete) Description() string {
	return setInCompletePrompt
}

func (s *SetInComplete) Run(storage *storage.Storage, attributes map[string]string, payload string) string {
	pos, _ := strconv.Atoi(payload)
	if storage.SetInComplete(pos) {
		return fmt.Sprintf("step %d marked as completed", pos)
	}
	return fmt.Sprintf("no plan step at position %d", pos)
}

func (s *SetInComplete) Timeout() *time.Duration {
	return nil
}

func (s *SetInComplete) ExamplePayload() *string {
	p := "2"
	return &p
}

func (s *SetInComplete) ExampleAttributes() map[string]string {
	return nil
}

func (s *SetInComplete) RequiredVariables() []*string {
	return nil
}

func (s *SetInComplete) RequiresUserConfirmation() bool {
	return true
}

func (s *SetInComplete) GetNamespace() types.NamespaceType {
	return types.PLANNING
}

func (s *SetInComplete) NamespaceDescription() string {
	return nsPrompt
}
