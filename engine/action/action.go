package action

import (
	"time"

	"github.com/runetale/notch/engine/events"
	"github.com/runetale/notch/types"
)

// all namespace's implement this interfaces
type Action interface {
	GetNamespace() types.NamespaceType
	Name() string
	Description() string
	// todo: chanelでstorageなどを送る
	Run(event events.Event, attributes map[string]string, payload string) string
	Timeout() time.Duration
	RequiredVariables() []string
	RequiresUserConfirmation() bool
}
