package state

import (
	"github.com/runetale/notch/llm"
)

type Execution struct {
	// llm response
	Response *string
	// parsed llm response to invocation
	invocation *llm.Invocation
	result     *string
	error      error
}

func NewExecution() *Execution {
	return &Execution{}
}
