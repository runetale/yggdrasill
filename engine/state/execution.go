// history of executed by engine
package state

import (
	"github.com/runetale/notch/llm"
)

type Execution struct {
	// llm response
	Response *string
	// parsed llm response to invocation
	Invocation *llm.Invocation
	Result     *string
	Error      error
}

func NewExecution() *Execution {
	return &Execution{}
}
