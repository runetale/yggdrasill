package tasklet

import (
	"os"
	"time"

	"github.com/runetale/notch/engine/action"
	"github.com/runetale/notch/engine/events"
	"github.com/runetale/notch/types"
)

type Tasklet struct {
	name             string
	description      string
	workingDirectory string
	maxShownOutput   uint32
	args             map[string]string
	examplePayload   *string
	timeout          string
	tool             string
}

func NewTasklet() action.Action {
	return &Tasklet{}
}

func (s *Tasklet) Name() string {
	return "shell"
}

func (s *Tasklet) Description() string {
	filepath := "shell.prompt"
	data, _ := os.ReadFile(filepath)
	return string(data)
}

func (s *Tasklet) Run(event events.Event, attributes map[string]string, payload string) string {
	return "run"
}

func (s *Tasklet) Timeout() time.Duration {
	return time.Duration(0)
}

func (s *Tasklet) RequiredVariables() []string {
	return []string{""}
}

func (s *Tasklet) RequiresUserConfirmation() bool {
	return true
}

func (s *Tasklet) GetNamespace() types.NamespaceType {
	return types.SHELL
}
