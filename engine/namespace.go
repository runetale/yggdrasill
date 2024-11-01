package engine

type Namespace string

const (
	FILESYSTEM Namespace = ""
	GOAL       Namespace = ""
	HTTP       Namespace = ""
	MEMORY     Namespace = ""
	SHELL      Namespace = ""
	TIME       Namespace = ""
	RAG        Namespace = ""
	PLANNING   Namespace = ""
	TASK       Namespace = ""
)
