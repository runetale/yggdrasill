package llm

type Invocation struct {
	action     string
	attributes map[string]string
	payload    string
}

func NewInvocation() *Invocation {
	return &Invocation{}
}
