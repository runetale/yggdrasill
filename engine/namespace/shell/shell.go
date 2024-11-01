package shell

import (
	"os"
	"time"

	"github.com/runetale/notch/types"
)

type Shell struct {
}

func NewShell() *Shell {
	return &Shell{}
}

func (s *Shell) Name() string {
	return "shell"
}

func (s *Shell) Description() string {
	filepath := "shell.prompt"
	data, _ := os.ReadFile(filepath)
	return string(data)
}

// func (s *Shell) Run(state state.State, attributes map[string]string, payload string) string {
// 	return "run"
// }

func (s *Shell) Timeout() time.Duration {
	return time.Duration(0)
}

func (s *Shell) RequiredVariables() []string {
	return []string{""}
}

func (s *Shell) RequiresUserConfirmation() bool {
	return true
}

func (s *Shell) GetNamespace() types.NamespaceType {
	return types.SHELL
}
