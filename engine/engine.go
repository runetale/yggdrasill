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
		event := <-e.channel.Receiver
		switch event.EventType() {
		case events.ActionExecuted:

		}
	}
}

func (e *Engine) automaton() {
	for {
		// prepare chat option
		option := e.prepareAutomaton()

		// request to chat
		invocations := []*llm.Invocation{}

		// NOTE:
		// 1.toolCalls
		// 2.response
		// 1,2 iteration creates an invocation
		toolCalls, response := e.factory.Chat(option)
		if toolCalls == nil {
			// use our strategy
			invocations = serializer.TryParse(&response)
		} else {
			// use native function call by model supports
			invocations = toolCalls
		}

		// チャネルに送信
		// on_state_update

		// llmからresponseとtool_callsをもらう
		// generator.chat

		// もらったtool_callsのinvocationsからを使用して、コマンドを実行
		// chat historyに追加
		// チャネルに送信

		// 再度実行無限ループ

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
