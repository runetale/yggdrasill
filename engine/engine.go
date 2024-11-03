package engine

import (
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
	task       *task.Tasklet

	waitCh chan struct{}
}

func NewEngine(t *task.Tasklet, c *llm.LLMFactory, maxIterations uint) *Engine {
	channel := events.NewChannel()

	serializationInvocationCb := func(inv *llm.Invocation) *string {
		return serializer.SerializeInvocation(inv)
	}

	s := state.NewState(channel, t, maxIterations, serializationInvocationCb)

	return &Engine{
		channel:    channel,
		factory:    c,
		maxHistory: t.GetMaxHistory(),
		state:      s,
		task:       t,
		waitCh:     make(chan struct{}),
	}
}

func (e *Engine) Start() {
	go e.consumeEvent()
	go e.automaton()
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
		// waiting receiver for each events
		event := <-e.channel.Chan
		switch event.EventType() {
		case events.ActionExecuted:

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
		invocations := []*llm.Invocation{}
		toolCalls, response := e.factory.Chat(option)
		if toolCalls == nil {
			// use our strategy
			invocations = serializer.TryParse(response)
		} else {
			// use native function call by model supports
			invocations = toolCalls
		}

		if invocations == nil {
			if response == "" {
				e.onEmptyResponse()
			} else {
				e.onInvalidResponse(response)
			}
		} else {
			e.onValidResponse()
		}

		// parsing invocations

		for _, inv := range invocations {
			e.state.GetAciton(inv.Action)

		}

		// Engineを修了する
		e.Stop()
	}
}

func (e *Engine) prepareAutomaton() *llm.ChatOption {
	e.state.OnEvent(events.NewEvent(events.MetricsUpdate, "engine", "prepare-automaton"))
	// get system prompt by state
	systemPrompt, err := serializer.DisplaySystemPrompt(e.state)
	if err != nil {

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
		sysp, err := serializer.DisplaySystemPrompt(e.state)
		if err != nil {
			panic(err)
		}
		options.UpdateSystemPrompt(sysp)

		// update history
		history := e.state.ToChatHistory(int(e.maxHistory))
		options.UpdateHistroy(history)
	}
	e.sendEvent(events.NewEvent(events.StateUpdate, "engine", "on-update-state"))
}

func (e *Engine) sendEvent(events *events.Event) {
	e.state.OnEvent(events)
}

func (e *Engine) onInvalidResponse(reponse string) {
	e.state.IncrementUnparsedMetrics()
	e.state.AddUnparsedResponseToHistory(reponse, "no effective solution found, follow the instructions to correct this")
	e.state.OnEvent(events.NewEvent(events.InvalidResponse, "engine", "on-invalid-response"))

}

func (e *Engine) onEmptyResponse() {
	e.state.IncrementEmptyMetrics()
	e.state.AddUnparsedResponseToHistory("", "return to empty response")
	e.state.OnEvent(events.NewEvent(events.EmptyResponse, "engine", "on-empty-response"))
}

func (e *Engine) onValidResponse() {
	e.state.IncrementValidMetrics()
}
