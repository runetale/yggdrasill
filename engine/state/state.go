// automaton machine state
package state

import (
	"fmt"
	"log"

	"github.com/runetale/notch/engine/action"
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
	complete chan bool

	// call from engine and storage
	onEventCallback func(event *events.Event)

	// serialize callback function
	SerializeInvocation func(inv *llm.Invocation) *string

	metrics *Metrics
}

// TODO implement rag model
func NewState(
	sender *events.Channel,
	task *task.Tasklet,
	maxIterations uint,
	serializationInvocation func(inv *llm.Invocation) *string,
) *State {
	namespaces := make([]*namespace.Namespace, 0)
	storages := make(map[string]*storage.Storage, 0)
	variables := make(map[string]string, 0)
	history := make([]*Execution, 0)

	s := &State{
		sender: sender,
	}
	s.SerializeInvocation = serializationInvocation

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
		s.sender.Chan <- event
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
	s.metrics = metrics

	return s
}

func (s *State) Complete() <-chan bool {
	return s.complete
}

func (s *State) Close() {
	s.complete <- true
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

// update history functions
func (s *State) AddUnparsedResponseToHistory(response string, err string) {
	s.history = append(s.history, NewExecution(&response, nil, nil, &err))
}

func (s *State) AddSuccessToHistory(invocation *llm.Invocation, result *string) {
	s.history = append(s.history, NewExecution(nil, invocation, result, nil))
}

func (s *State) AddErrorToHistory(invocation *llm.Invocation, err string) {
	s.history = append(s.history, NewExecution(nil, invocation, nil, &err))
}

// when this function called from `first chat“ and `on state update`
func (s *State) ToChatHistory(max int) []*llm.Message {
	var latest []*Execution
	if len(s.history) > max {
		latest = s.history[:max+1]
	} else {
		latest = s.history
	}

	if latest == nil {
		return nil
	}

	// to messages
	history := []*llm.Message{}
	for _, entry := range latest {
		// agent messages
		if entry.Response != nil {
			history = append(history, &llm.Message{
				MessageType: llm.AGETNT,
				Response:    entry.Response,
				Invocation:  nil,
			})
		} else if entry.Invocation != nil {
			history = append(history, &llm.Message{
				MessageType: llm.AGETNT,
				// parse to invocation to string,
				// to including the results of executing a "function call" when executing factory.Chat()
				Response:   s.SerializeInvocation(entry.Invocation),
				Invocation: entry.Invocation,
			})
		}

		// feedback messages
		var res string
		if entry.Error != nil {
			res = fmt.Sprintf("ERROR: %s", *entry.Error)
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

// metrics functions
func (s *State) IncrementEmptyMetrics() {
	s.metrics.errors.emptyResponses += 1
}

func (s *State) IncrementUnparsedMetrics() {
	s.metrics.errors.unparsedResonses += 1
}

func (s *State) IncrementUnknownMetrics() {
	s.metrics.errors.unknownActions += 1
}

func (s *State) IncrementValidMetrics() {
	s.metrics.validResponses += 1
}

func (s *State) GetAciton(actionName string) action.Action {
	for _, group := range s.namespaces {
		for _, ac := range group.GetActions() {
			if actionName == ac.Name() {
				return ac
			}
		}
	}
	return nil
}
