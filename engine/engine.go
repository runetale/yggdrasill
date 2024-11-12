package engine

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/runetale/notch/engine/action"
	"github.com/runetale/notch/engine/chat"
	"github.com/runetale/notch/engine/serializer"
	"github.com/runetale/notch/engine/state"
	"github.com/runetale/notch/events"
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
	saveTo     *string

	waitCh chan struct{}
}

func NewEngine(t *task.Task, c *llm.LLMFactory, maxIterations uint, nativeTool bool, saveTo *string) *Engine {
	channel := events.NewChannel()

	serializationInvocationCb := func(inv *chat.Invocation) *string {
		return serializer.SerializeInvocation(inv)
	}
	s := state.NewState(channel, t, maxIterations, serializationInvocationCb)

	return &Engine{
		channel:    channel,
		factory:    c,
		maxHistory: t.GetMaxHistory(),
		state:      s,
		task:       t,
		timeout:    s.GetTask().GetTimeout(),
		nativeTool: nativeTool,
		saveTo:     saveTo,

		waitCh: make(chan struct{}),
	}
}

func (e *Engine) Start() {
	go e.consumeEvent()
	go e.automaton()
}

func (e *Engine) Stop() {
	log.Printf("shutdown...")
	close(e.waitCh)
}

func (e *Engine) Done() <-chan struct{} {
	return e.waitCh
}

// for only display
func (e *Engine) consumeEvent() {
	for {
		// waiting event chan for each events
		event := <-e.channel.Chan
		log.Println(event.Display())
	}
}

func (e *Engine) automaton() {
	for {
		// prepare chat option
		option := e.prepareAutomaton()

		// update state event
		e.OnUpdateState(option, false)

		// response from llm
		var invocations []*chat.Invocation
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
				e.onInvalidAction(inv, nil)
				break
			}

			// validate actions
			err := inv.ValidateAction(ac)
			if err != nil {
				errStr := err.Error()
				e.onInvalidAction(inv, &errStr)
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
	fmt.Printf("run for %s\n", ac.Name())
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

func (e *Engine) prepareAutomaton() *chat.ChatOption {
	e.state.OnEvent(events.NewMetricsEvent(e.state.DisplayMetrics()))
	// get system prompt by state
	systemPrompt, err := serializer.DisplaySystemPrompt(e.state)
	if err != nil {
		log.Fatalf("prepare automaton error %s", err.Error())
	}

	// get prompt by state
	prompt := e.task.GetPrompt()

	// to chat history, return the history of messsagesb by maxHistory count
	history := e.state.ToChatHistory(int(e.maxHistory))

	return chat.NewChatOption(systemPrompt, prompt, history)
}

func (e *Engine) OnUpdateState(options *chat.ChatOption, refresh bool) {
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
	stringsList := make([]string, len(options.GetHistory()))
	for i, h := range options.GetHistory() {
		stringsList[i] = h.Display()
	}
	e.state.OnEvent(events.NewStateUpdateEvent(options.GetSystemPrompt(), options.GetPrompt(), strings.Join(stringsList, "\n"), e.saveTo))
}

func (e *Engine) onEmptyResponse() {
	e.state.IncrementEmptyMetrics()
	e.state.AddUnparsedResponseToHistory("", "return to empty response")
	e.state.OnEvent(events.NewEmptyResponseEvent())
}

func (e *Engine) onInvalidResponse(response string) {
	e.state.IncrementUnparsedMetrics()
	e.state.AddUnparsedResponseToHistory(response, "no effective solution found, follow the instructions to correct this")
	e.state.OnEvent(events.NewInvalidResponseEvent(response))
}

func (e *Engine) onValidResponse() {
	e.state.IncrementValidMetrics()
}

func (e *Engine) onValidAction() {
	e.state.IncrementValidActionsMetrics()
}

func (e *Engine) onInvalidAction(inv *chat.Invocation, err *string) {
	e.state.IncrementUnknownMetrics()
	e.state.AddErrorToHistory(inv, err)
	e.state.OnEvent(events.NewInvalidActionEvent(inv.Action, *err))
}

func (e *Engine) onTimeoutAction(inv *chat.Invocation, start time.Duration) {
	e.state.IncrementTimeoutActionMetrics()
	err := "action time out"
	e.state.AddErrorToHistory(inv, &err)
	e.state.OnEvent(events.NewActionTimeoutEvent(inv.Action, start))
}

func (e *Engine) onExecutedErrorAction(inv *chat.Invocation, err *string, start time.Duration) {
	e.state.IncrementErroredActionMetrics()
	e.state.AddErrorToHistory(inv, err)
	e.state.OnEvent(events.NewActionExecutedEvent(inv.Action, err, nil, start))
}

func (e *Engine) onExecutedSuccessAction(inv *chat.Invocation, result *string, start time.Duration) {
	e.state.IncrementSuccessActionMetrics()
	e.state.AddSuccessToHistory(inv, result)
	e.state.OnEvent(events.NewActionExecutedEvent(inv.Action, nil, result, start))
}
