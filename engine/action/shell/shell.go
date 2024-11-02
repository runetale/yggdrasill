package shell

import (
	"os"
	"time"

	"github.com/runetale/notch/engine/action"
	"github.com/runetale/notch/storage"
	"github.com/runetale/notch/types"
)

type Shell struct {
	storageType types.StorageType
	predefined  *map[string]string
}

func NewShell() action.Action {
	return &Shell{
		storageType: types.UNTAGGED,
		predefined:  nil,
	}
}

func (s *Shell) StorageType() types.StorageType {
	return s.storageType
}

func (s *Shell) Predefined() *map[string]string {
	return s.predefined
}

func (s *Shell) Name() string {
	return "shell"
}

func (s *Shell) Description() string {
	filepath := "shell.prompt"
	data, _ := os.ReadFile(filepath)
	return string(data)
}

func (s *Shell) Run(storage storage.Storage, attributes map[string]string, payload string) string {
	return "run"
}

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
