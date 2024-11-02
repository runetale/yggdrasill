package action

import (
	"time"

	"github.com/runetale/notch/storage"
	"github.com/runetale/notch/types"
)

// all namespace's implement this interfaces
type Action interface {
	GetNamespace() types.NamespaceType
	Name() string
	StorageType() types.StorageType
	Predefined() *map[string]string
	Description() string
	Run(event storage.Storage, attributes map[string]string, payload string) string
	Timeout() time.Duration
	// "$SSH_HOST"などのvariablesを設定した場合に取得する
	RequiredVariables() []string
	RequiresUserConfirmation() bool
}
