package types

type NamespaceType string

// TODO: add more namespaces of actions
// - taking screenshot
// - moving mouse(robotogo)
// - ui interactions

const (
	FILESYSTEM NamespaceType = "filesytem"
	GOAL       NamespaceType = "goal"
	HTTP       NamespaceType = "http"
	MEMORY     NamespaceType = "memory"
	SHELL      NamespaceType = "shell"
	TIME       NamespaceType = "time"
	RAG        NamespaceType = "rag"
	PLANNING   NamespaceType = "planning"
	TASK       NamespaceType = "task"
	CUSTOM     NamespaceType = "custom"
)

func GetNameSpaceValues() []NamespaceType {
	ns := []NamespaceType{}
	ns = append(ns, FILESYSTEM)
	ns = append(ns, GOAL)
	ns = append(ns, HTTP)
	ns = append(ns, MEMORY)
	ns = append(ns, SHELL)
	ns = append(ns, TIME)
	ns = append(ns, RAG)
	ns = append(ns, PLANNING)
	ns = append(ns, TASK)
	ns = append(ns, CUSTOM)
	return ns
}
