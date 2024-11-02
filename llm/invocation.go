package llm

type Invocation struct {
	Action     string
	Attributes map[string]string
	Payload    string
}

func NewInvocation() *Invocation {
	return &Invocation{}
}
