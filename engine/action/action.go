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
	Description() string
	Run(event storage.Storage, attributes map[string]string, payload string) string
	Timeout() time.Duration
	RequiredVariables() []string
	RequiresUserConfirmation() bool
}
