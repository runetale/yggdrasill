package engine

import (
	"github.com/runetale/notch/engine/events"
	"github.com/runetale/notch/engine/state"
	"github.com/runetale/notch/llm"
	"github.com/runetale/notch/task"
)

type Engine struct {
	channel *events.Channel
	client  *llm.LLMClient
	state   *state.State
	task    *task.Tasklet

	closeCh chan struct{}
	waitCh  chan struct{}
}

func NewEngine(t *task.Tasklet, c *llm.LLMClient) *Engine {
	channel := events.NewChannel()

	s := state.NewState(channel, t, 0)

	return &Engine{
		channel: channel,
		client:  c,
		state:   s,
		task:    t,
	}
}

func (e *Engine) Start() {
	go e.consumeEvent()
	go e.automaton()
}

func (e *Engine) Stop() {
	<-e.closeCh
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
	// chat historyを生成
	e.prepareAutomaton()

	// チャネルに送信
	// on_state_update

	// llmからresponseとtool_callsをもらう
	// generator.chat

	// もらったtool_callsのinvocationsからを使用して、コマンドを実行
	// chat historyに追加
	// チャネルに送信

	// 再度実行無限ループ

	// Engineを修了する
	e.closeCh <- struct{}{}
}

func (e *Engine) prepareAutomaton() *llm.ChatOption {
	// get system prompt by state
	e.state.OnEvent(events.NewEvent(events.MetricsUpdate, "engine", "prepare-automaton"))

	// get prompt by state

	// get history by state
	return llm.NewChatOption("", "", nil)
}
