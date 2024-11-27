package action

import (
	"time"

	"github.com/runetale/yggdrasill/storage"
	"github.com/runetale/yggdrasill/types"
)

// all namespace's implement this interfaces
type Action interface {
	GetNamespace() types.NamespaceType
	Name() string
	// namespace description
	NamespaceDescription() string
	// namespace action description
	Description() string
	Run(s *storage.Storage, attributes map[string]string, payload string) string
	Timeout() *time.Duration
	RequiredVariables() []*string // retrieved when variables such as `$SSH_HOST` are set
	RequiresUserConfirmation() bool
	ExamplePayload() *string
	ExampleAttributes() map[string]string
}
