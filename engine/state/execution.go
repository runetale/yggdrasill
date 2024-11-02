package state

import (
	"github.com/runetale/notch/llm"
)

type Execution struct {
	// unparsed response caused an error
	Response *string
	// parsed invocation
	Invocation *llm.Invocation
	result     *string
	error      error
}

func NewExecution() *Execution {
	return &Execution{}
}
