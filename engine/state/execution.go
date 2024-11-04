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
	// if engine executed success
	Result *string
	// if engine executed error
	Error *string
}

func NewExecution(
	response *string,
	invocation *llm.Invocation,
	result *string,
	err *string,
) *Execution {
	return &Execution{
		Response:   response,
		Invocation: invocation,
		Result:     result,
		Error:      err,
	}
}
