// automaton machine state
package state

import (
	"fmt"
	"log"

	"github.com/runetale/notch/engine/events"
	"github.com/runetale/notch/engine/namespace"
	"github.com/runetale/notch/storage"
	"github.com/runetale/notch/task"
	"github.com/runetale/notch/types"
)

type State struct {
	Task      task.Tasklet
	Storages  map[string]storage.Storage
	Variables map[string]string // pre-define variables
	// 各Namespaceを持った構造体
	Namespaces []*namespace.Namespace
	History    []*Execution // execed histories

	// sent to engine.consumeEvent
	sender   chan<- events.Channel
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
			namespaces = append(namespaces, namespace.NewNamespace(o, nil))
		}
	} else {
		// adding only task defined namespaces
		for _, o := range using {
			if *o == "*" {
				*o = ""
				continue
			}
			namespaces = append(namespaces, namespace.NewNamespace(types.NamespaceType(*o), nil))
		}
	}

	// set variables
	for _, o := range namespaces {
		required := o.Action.RequiredVariables()
		log.Printf("actions %s requires %v\n", o.Action.Name(), required)
		for _, vn := range required {
			exp := fmt.Sprintf("$%s", vn)
			varname, value, err := task.ParseVariableExpr(exp)
			if err != nil {
				return nil
			}
			variables[varname] = value
		}
	}

	// add task defined actions by yaml
	functions := task.GetFunctions()
	namespaces = append(namespaces, namespace.NewNamespace(types.CUSTOM, functions))

	// create storages by namespaces

	return &State{
		Namespaces: namespaces,
		Variables:  variables,
		History:    history,
		sender:     sender,
	}
}
