package planning

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/runetale/notch/engine/action"
	"github.com/runetale/notch/storage"
	"github.com/runetale/notch/types"
)

type SetComplete struct {
}

func NewSetComplete() action.Action {
	return &SetComplete{}
}

func (s *SetComplete) Name() string {
	return "set_step_completed"
}

func (s *SetComplete) Description() string {
	filepath := "set-complete.prompt"
	data, _ := os.ReadFile(filepath)
	return string(data)
}

func (s *SetComplete) Run(storage *storage.Storage, attributes map[string]string, payload string) string {
	pos, _ := strconv.Atoi(payload)
	if storage.SetComplete(pos) {
		return fmt.Sprintf("step %d marked as completed", pos)
	}
	return fmt.Sprintf("no plan step at position %d", pos)
}

func (s *SetComplete) Timeout() *time.Duration {
	return nil
}

func (s *SetComplete) ExamplePayload() *string {
	p := "2"
	return &p
}

func (s *SetComplete) ExampleAttributes() map[string]string {
	return nil
}

func (s *SetComplete) RequiredVariables() []*string {
	return nil
}

func (s *SetComplete) RequiresUserConfirmation() bool {
	return true
}

func (s *SetComplete) GetNamespace() types.NamespaceType {
	return types.PLANNING
}

func (s *SetComplete) NamespaceDescription() string {
	filepath := "ns.prompt"
	data, _ := os.ReadFile(filepath)
	return string(data)
}
