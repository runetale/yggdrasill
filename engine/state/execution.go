package state

type Execution struct {
	// unparsed response caused an error
	Response *string
	// parsed invocation
	Invocation *Invocation
	result     *string
	error      error
}

func NewExecution() *Execution {
	return &Execution{}
}
