// automaton machine state
package state

import (
	"github.com/runetale/notch/engine/events"
	"github.com/runetale/notch/engine/namespace"
	"github.com/runetale/notch/storage"
	"github.com/runetale/notch/task"
	"github.com/runetale/notch/types"
)

type State struct {
	Task      task.Tasklet
	Storages  map[string]storage.Storage
	Variables map[string]string // pre-define
	// 各Namespaceを持った構造体
	Namespaces []*namespace.Namespace
	History    []*Execution // execed histories

	sender   events.Channel
	complete bool
}

// TODO implement rag model
func NewState(
	sender chan<- events.Channel,
	task *task.Tasklet,
	maxIterations uint64,
) *State {
	namespaces := make([]*namespace.Namespace, 0)
	variables := make(map[string]string, 0)
	history := make([]*Execution, 0)

	// get namespaces
	using := task.GetUsing()
	if len(using) == 0 {
		// creating default namespaces
		ns := types.GetNameSpaceValues()
		for _, o := range ns {
			namespaces = append(namespaces, namespace.NewNamespace(o))
		}
	} else {
		// adding only task defined namespaces
		for _, o := range using {
			if *o == "*" {
				*o = ""
				continue
			}
			namespaces = append(namespaces, namespace.NewNamespace(types.NamespaceType(*o)))
		}
	}

	// set variables
	// setされたすべてのnamespaceの値を呼び出す
	for _, o := range namespaces {
	}

	return &State{
		Namespaces: namespaces,
		Variables:  variables,
		History:    history,
	}
}
