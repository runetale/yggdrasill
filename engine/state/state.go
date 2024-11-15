// automaton machine state
package state

import (
	"fmt"
	"log"

	"github.com/runetale/notch/engine/action"
	"github.com/runetale/notch/engine/chat"
	"github.com/runetale/notch/engine/namespace"
	"github.com/runetale/notch/events"
	"github.com/runetale/notch/storage"
	"github.com/runetale/notch/task"
	"github.com/runetale/notch/types"
)

type State struct {
	task       *task.Task
	storages   map[string]*storage.Storage
	variables  map[string]string      // pre-define variables
	namespaces []*namespace.Namespace // user defined each namespaces
	history    []*Execution           // execed histories

	// sent to engine.consumeEvent
	sender *events.Channel

	// call from engine and storage
	onEventCallback func(event events.DisplayEvent)

	// serialize callback function
	SerializeInvocation func(inv *chat.Invocation) *string

	metrics *Metrics
}

// TODO implement rag model
func NewState(
	sender *events.Channel,
	task *task.Task,
	maxIterations uint,
	serializationInvocation func(inv *chat.Invocation) *string,
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
		for _, action := range o.Actions() {
			required := action.RequiredVariables()
			if required == nil {
				break
			}
			log.Printf("actions %s requires %v\n", action.Name(), required)
			for _, vn := range required {
				exp := fmt.Sprintf("$%s", *vn)
				varname, value, err := task.ParseVariableExpr(exp)
				if err != nil {
					log.Fatalf("error parse variable expr %s", err.Error())
					return nil
				}
				variables[varname] = value
			}
		}
	}

	// TODO: check the custom functions
	// add task defined actions by yaml, if user's was set
	if task.GetFunctions() != nil {
		functions := task.GetFunctions()
		namespaces = append(namespaces, namespace.NewNamespace(types.NamespaceType(task.GetName()), functions))
	}

	// set callback function
	onEventCallback := func(event events.DisplayEvent) {
		go func() {
			s.sender.Chan <- event
		}()
	}
	s.onEventCallback = onEventCallback

	// create storages by namespaces, if storages descriptor have it by each namespaces
	for _, namespace := range namespaces {
		if namespace.GetStrorageDescriptor() != nil {
			for _, storageDescriptor := range namespace.GetStrorageDescriptor() {
				if _, exists := storages[storageDescriptor.Name()]; !exists {
					newStorage := storage.NewStorage(storageDescriptor.Name(), storageDescriptor.Type(), onEventCallback)

					if storageDescriptor.Predefined() != nil {
						for key, value := range storageDescriptor.Predefined() {
							newStorage.AddData(key, *value)
						}
					}
					log.Printf("create storage [%s]\n", storageDescriptor.Name())
					storages[storageDescriptor.Name()] = newStorage
				}
			}
		}
	}

	// new state
	s.task = task
	s.namespaces = namespaces
	s.variables = variables
	s.history = history
	s.storages = storages

	// if the goal namespace is enabled, set the current goal
	for key, s := range storages {
		if key == "goal" {
			prompt := task.GetPrompt()
			log.Printf("set goal prompt => '%s'\n", prompt)
			s.SetCurrent(prompt)
		}
	}

	// set metrics
	metrics := NewMetrics(uint(maxIterations))
	s.metrics = metrics

	return s
}

// called from engine
func (s *State) OnEvent(event events.DisplayEvent) {
	s.onEventCallback(event)
}

func (s *State) GetTask() *task.Task {
	return s.task
}

func (s *State) GetPrompt() string {
	return s.task.GetPrompt()
}

func (s *State) GetStorages() map[string]*storage.Storage {
	return s.storages
}

func (s *State) GetStorage(actionName string) *storage.Storage {
	return s.storages[actionName]
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

func (s *State) AddSuccessToHistory(invocation *chat.Invocation, result *string) {
	s.history = append(s.history, NewExecution(nil, invocation, result, nil))
}

func (s *State) AddErrorToHistory(invocation *chat.Invocation, err *string) {
	s.history = append(s.history, NewExecution(nil, invocation, nil, err))
}

// when this function called from `first chat“ and `on state update`
func (s *State) ToChatHistory(max int) []*chat.Message {
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
	// todo: historyの内容が正しいか？
	history := []*chat.Message{}
	for _, entry := range latest {
		// agent messages
		if entry.Response != nil {
			history = append(history, &chat.Message{
				MessageType: chat.AGETNT,
				Response:    entry.Response,
				Invocation:  nil,
			})
		} else if entry.Invocation != nil {
			history = append(history, &chat.Message{
				MessageType: chat.AGETNT,
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

		history = append(history, &chat.Message{
			MessageType: chat.FEEDBACK,
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
	s.metrics.errors.unparsedResponses += 1
}

func (s *State) IncrementUnknownMetrics() {
	s.metrics.errors.unknownActions += 1
}

func (s *State) IncrementValidMetrics() {
	s.metrics.validResponses += 1
}

func (s *State) IncrementValidActionsMetrics() {
	s.metrics.validActions += 1
}

func (s *State) IncrementErroredActionMetrics() {
	s.metrics.errors.erroredActions += 1
}

func (s *State) IncrementSuccessActionMetrics() {
	s.metrics.successActions += 1
}

func (s *State) IncrementTimeoutActionMetrics() {
	s.metrics.errors.timedoutActions += 1
}

// calling invocation.action from engine
// get the namespace of the specified action name
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

func (s *State) DisplayMetrics() string {
	return s.metrics.Display()
}
