package action

import (
	"time"

	"github.com/runetale/notch/engine/state"
	"github.com/runetale/notch/types"
)

// all namespace's implement this interfaces
// このinterfaceのカプセル化をうまく行う
type Action interface {
	GetNamespace() types.NamespaceType
	Name() string
	Description() string
	// これをどう実行するか考える
	Run(state state.State, attributes map[string]string, payload string) string
	Timeout() time.Duration
	RequiredVariables() []string
	RequiresUserConfirmation() bool
}
