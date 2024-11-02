// automaton machine state
package state

import (
	"fmt"
	"log"

	"github.com/runetale/notch/engine/events"
	"github.com/runetale/notch/engine/namespace"
	"github.com/runetale/notch/llm"
	"github.com/runetale/notch/storage"
	"github.com/runetale/notch/task"
	"github.com/runetale/notch/types"
)

type State struct {
	task      *task.Tasklet
	storages  map[string]*storage.Storage
	variables map[string]string // pre-define variables
	// 各Namespaceを持った構造体
	namespaces []*namespace.Namespace
	history    []*Execution // execed histories

	// sent to engine.consumeEvent
	sender   *events.Channel
	complete bool

	// call from engine and storage
	onEventCallback func(event *events.Event)

	metrics *Metrics
}

// TODO implement rag model
func NewState(
	sender *events.Channel,
	task *task.Tasklet,
	maxIterations uint,
) *State {
	namespaces := make([]*namespace.Namespace, 0)
	storages := make(map[string]*storage.Storage, 0)
	variables := make(map[string]string, 0)
	history := make([]*Execution, 0)
	complete := false
	s := &State{
		sender: sender,
	}

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
		required := o.GetAction().RequiredVariables()
		log.Printf("actions %s requires %v\n", o.GetAction().Name(), required)
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

	// set callback function
	onEventCallback := func(event *events.Event) {
		s.sender.Sender <- event
	}
	s.onEventCallback = onEventCallback

	// create storages by namespaces
	for _, namespace := range namespaces {
		for _, currentStorage := range namespace.GetStorages() {
			// if storages nil, set to newstorage by namespace
			if currentStorage == nil {
				newStorage := storage.NewStorage(namespace.GetName(), namespace.GetStorageType(), onEventCallback)
				// set namespace to storage
				currentStorage = newStorage
				storages[namespace.GetName()] = newStorage
			}
		}
	}
	for key, storage := range storages {
		if key == "goal" {
			prompt := task.GetPrompt()
			fmt.Println("set goal prompt")
			fmt.Printf("%s\n", prompt)
			storage.SetCurrent(prompt)
		}
	}

	metrics := NewMetrics(uint(maxIterations))

	s.task = task
	s.storages = storages
	s.namespaces = namespaces
	s.variables = variables
	s.history = history
	s.sender = sender
	s.complete = complete
	s.metrics = metrics

	return s
}

// called from engine
func (s *State) OnEvent(event *events.Event) {
	s.onEventCallback(event)
}

func (s *State) GetTask() *task.Tasklet {
	return s.task
}

func (s *State) GetPrompt() string {
	if s.task.Prompt != nil {
		return *s.task.Prompt
	}
	return "state prompt not set"
}

func (s *State) GetStorages() map[string]*storage.Storage {
	return s.storages
}

func (s *State) GetNamespaces() []*namespace.Namespace {
	return s.namespaces
}

func (s *State) GetMaxIteration() uint {
	return s.metrics.maxStep
}

func (s *State) GetCurrentStep() uint {
	return s.metrics.currentStep
}

func (s *State) GetChatHistory(max int) []*llm.Message {
	var latest []*Execution
	if len(s.history) > max {
		latest = s.history[:max+1]
	} else {
		latest = s.history
	}

	history := []*llm.Message{}
	for _, entry := range latest {
		if entry.Response != nil {
			history = append(history, &llm.Message{
				MessageType: llm.AGETNT,
				Response:    entry.Response,
				Invocation:  nil,
			})
		} else if entry.Invocation != nil {
			history = append(history, &llm.Message{
				MessageType: llm.AGETNT,
				Response:    nil,
				Invocation:  entry.Invocation,
			})
		}

		var res string
		if entry.Error != nil {
			res = fmt.Sprintf("ERROR: %s", entry.Error.Error())
		} else if entry.Result != nil {
			res = *entry.Result
		} else {
			res = ""
		}

		history = append(history, &llm.Message{
			MessageType: llm.FEEDBACK,
			Response:    &res,
			Invocation:  entry.Invocation,
		})
	}

	return history
}
