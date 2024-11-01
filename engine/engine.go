package engine

import (
	"github.com/runetale/notch/engine/events"
	"github.com/runetale/notch/engine/state"
	"github.com/runetale/notch/llm"
	"github.com/runetale/notch/task"
)

type Engine struct {
	// stateの更新を行う
	Channel *events.Channel
	Client  llm.LLMClinet
	State   state.State
	Task    task.Tasklet

	closeCh chan struct{}
	waitCh  chan struct{}
}

func NewEngine() *Engine {
	channel := events.NewChannel()
	return &Engine{
		Channel: channel,
	}
}

func (e *Engine) Start() {
	go e.consumeEvent()

	go e.Automaton()
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
		event := <-e.Channel.Receiver
		switch event {
		case events.MetricsUpdate:

		}
	}
}

func (e *Engine) Automaton() {
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

func (e *Engine) prepareAutomaton() llm.ChatOption {
	// get system prompt by state

	// get prompt by state

	// get history by state
	return *llm.NewChatOption("", "", []string{""})
}
