package engine

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/runetale/notch/engine/action"
	"github.com/runetale/notch/engine/events"
	"github.com/runetale/notch/engine/serializer"
	"github.com/runetale/notch/engine/state"
	"github.com/runetale/notch/llm"
	"github.com/runetale/notch/task"
)

type Engine struct {
	channel    *events.Channel
	factory    *llm.LLMFactory
	state      *state.State
	maxHistory uint
	task       *task.Task
	timeout    *time.Duration
	nativeTool bool

	waitCh chan struct{}
}

func NewEngine(t *task.Task, c *llm.LLMFactory, maxIterations uint, nativeTool bool) *Engine {
	channel := events.NewChannel()

	serializationInvocationCb := func(inv *llm.Invocation) *string {
		return serializer.SerializeInvocation(inv)
	}
	s := state.NewState(channel, t, maxIterations, serializationInvocationCb)

	// check using native tools

	return &Engine{
		channel:    channel,
		factory:    c,
		maxHistory: t.GetMaxHistory(),
		state:      s,
		task:       t,
		timeout:    s.GetTask().GetTimeout(),
		nativeTool: nativeTool,
		waitCh:     make(chan struct{}),
	}
}

func (e *Engine) Start() {
	go e.consumeEvent()
	go e.automaton()

	// waiting terminated engine process
	comp := <-e.state.Complete()
	if comp {
		log.Printf("shutdown...")
		e.Stop()
	}
}

func (e *Engine) Stop() {
	close(e.waitCh)
}

func (e *Engine) Done() <-chan struct{} {
	return e.waitCh
}

// for only display
func (e *Engine) consumeEvent() {
	for {
		// waiting event cahn for each events
		event := <-e.channel.Chan
		switch event.EventType() {
		case events.MetricsUpdate:
		case events.StorageUpdate:
		case events.StateUpdate:
		case events.InvalidUpdate:
		case events.InvalidAction:
		case events.InvalidResponse:
		case events.ActionTimeOut:
		case events.ActionExecuted:
		case events.TaskComplete:
		case events.EmptyResponse:

		}
	}
}

func (e *Engine) automaton() {
	for {
		// prepare chat option
		option := e.prepareAutomaton()

		// update state event
		e.OnUpdateState(option, false)

		// response from llm
		var invocations []*llm.Invocation
		toolCalls, response := e.factory.Chat(option, e.nativeTool, e.state.GetNamespaces())

		// use our strategy
		if len(toolCalls) == 0 {
			invocations = serializer.TryParse(response)
		} else {
			// use native function call by model supports
			invocations = toolCalls
		}

		// return to llm response was null
		if len(invocations) == 0 {
			if response == "" {
				e.onEmptyResponse()
				continue
			} else {
				e.onInvalidResponse(response)
				continue
			}
		}

		// update metrics
		e.onValidResponse()

		// parsing invocations
		for _, inv := range invocations {
			// found action
			ac := e.state.GetAciton(inv.Action)
			if ac == nil {
				e.onInvalidAction(inv, fmt.Sprintf("cannot found action by %s", inv.Action))
				break
			}

			// validate actions
			err := inv.ValidateAction(ac)
			if err != nil {
				e.onInvalidAction(inv, fmt.Sprintf("invalid action %s", err.Error()))
				break
			}

			// update metrics
			e.onValidAction()

			// timeout
			timout := e.GetTimeout(ac)

			exec := true

			// y or n
			if ac.RequiresUserConfirmation() {
				log.Println("Warning: user confirmation required")
				start := time.Now()
				inp := "nope"

				for inp != "" && inp != "n" && inp != "y" {
					log.Println("invocation by y or n")
					inp = e.task.GetUserInput(fmt.Sprintf("%s [Yn] ", inv.FunctionCallString()))
					inp = strings.ToLower(inp)
				}

				if inp == "n" {
					log.Println("Warning: invocation rejected by user")
					elapsed := time.Since(start)
					h := fmt.Sprintf("invocation rejected. Elapsed time: %v\n", elapsed)
					e.onExecutedErrorAction(inv, &h, elapsed)
					exec = false
				}
			}

			// exec
			if exec {
				start := time.Now()
				result, err := e.timeoutRun(ac, timout, inv.Attributes, *inv.Payload)
				if err != nil {
					e.onTimeoutAction(inv, time.Since(start))
				}
				e.onExecutedSuccessAction(inv, &result, time.Since(start))
			}
		}

		// update state
		e.OnUpdateState(option, true)
		continue
	}
}

func (e *Engine) timeoutRun(ac action.Action, timeout time.Duration, attributes map[string]string, payload string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	result := make(chan string, 1)
	go func() {
		result <- ac.Run(e.state.GetStorage(ac.Name()), attributes, payload)
	}()

	select {
	case ret := <-result:
		return ret, nil
	case <-ctx.Done():
		return "", errors.New("operation timed out")
	}
}

func (e *Engine) GetTimeout(ac action.Action) time.Duration {
	var defaultTimeout = time.Hour * 24 * 30 // デフォルトは1ヶ月 (30日)
	if actionTimeout := ac.Timeout(); actionTimeout != nil {
		return *actionTimeout
	} else if e.timeout != nil {
		return *e.timeout
	}
	return defaultTimeout
}

func (e *Engine) prepareAutomaton() *llm.ChatOption {
	e.state.OnEvent(events.NewEvent(events.MetricsUpdate, "engine", "prepare-automaton"))
	// get system prompt by state
	systemPrompt, err := serializer.DisplaySystemPrompt(e.state)
	if err != nil {
		log.Fatalf("prepare automaton error %s", err.Error())
	}

	// get prompt by state
	prompt := e.task.GetPrompt()

	// to chat history, return the history of messsagesb by maxHistory count
	history := e.state.ToChatHistory(int(e.maxHistory))

	return llm.NewChatOption(systemPrompt, prompt, history)
}

func (e *Engine) OnUpdateState(options *llm.ChatOption, refresh bool) {
	if refresh {
		// update prompt
		sysprompt, err := serializer.DisplaySystemPrompt(e.state)
		if err != nil {
			log.Printf("error on update state %s", err.Error())
		}
		options.UpdateSystemPrompt(sysprompt)

		// update history
		history := e.state.ToChatHistory(int(e.maxHistory))
		options.UpdateHistroy(history)
	}
	e.sendEvent(events.NewEvent(events.StateUpdate, "engine", "on-update-state"))
}

func (e *Engine) sendEvent(events *events.Event) {
	e.state.OnEvent(events)
}

func (e *Engine) onInvalidResponse(response string) {
	e.state.IncrementUnparsedMetrics()
	e.state.AddUnparsedResponseToHistory(response, "no effective solution found, follow the instructions to correct this")
	e.state.OnEvent(events.NewEvent(events.InvalidResponse, "engine", fmt.Sprintf("agent did not provide valid instructions: \n\n%s\n\n", response)))
}

func (e *Engine) onEmptyResponse() {
	e.state.IncrementEmptyMetrics()
	e.state.AddUnparsedResponseToHistory("", "return to empty response")
	e.state.OnEvent(events.NewEvent(events.EmptyResponse, "engine", "on-empty-response"))
}

func (e *Engine) onValidResponse() {
	e.state.IncrementValidMetrics()
}

func (e *Engine) onValidAction() {
	e.state.IncrementValidActionsMetrics()
}

func (e *Engine) onInvalidAction(inv *llm.Invocation, err string) {
	e.state.IncrementUnknownMetrics()
	e.state.AddErrorToHistory(inv, err)
	e.state.OnEvent(events.NewEvent(events.InvalidAction, "engine", fmt.Sprintf("on-invalid-action %s", err)))
}

func (e *Engine) onExecutedErrorAction(inv *llm.Invocation, err *string, start time.Duration) {
	e.state.IncrementErroredActionMetrics()
	e.state.AddErrorToHistory(inv, *err)
	e.state.OnEvent(events.NewEvent(events.ActionExecuted, "engine", "on-executed-error-action"))
}

func (e *Engine) onExecutedSuccessAction(inv *llm.Invocation, result *string, start time.Duration) {
	e.state.IncrementSuccessActionMetrics()
	e.state.AddSuccessToHistory(inv, result)
	e.state.OnEvent(events.NewEvent(events.ActionExecuted, "engine", "on-executed-success-action"))
}

func (e *Engine) onTimeoutAction(inv *llm.Invocation, start time.Duration) {
	e.state.IncrementTimeoutActionMetrics()
	e.state.AddErrorToHistory(inv, "action time out")
	e.state.OnEvent(events.NewEvent(events.ActionTimeOut, "engine", "on-timeout-action"))
}
