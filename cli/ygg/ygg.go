package ygg

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/peterbourgon/ff/v2/ffcli"
	"github.com/runetale/yggdrasill/engine"
	"github.com/runetale/yggdrasill/llm"
	"github.com/runetale/yggdrasill/task"
)

const version = "0.0.1"

var yggArgs struct {
	taskpath      string
	prompt        string
	generator     string
	contextWindow int
	apiKey        string
	maxIterations int
	strategy      string
	forceFormat   bool
	saveTo        string
}

type StrategyFormat string

const (
	XML StrategyFormat = "xml"
)

var YggCmd = &ffcli.Command{
	Name:       "up",
	ShortUsage: "up [flags]",
	ShortHelp:  "command to start ygg",
	FlagSet: (func() *flag.FlagSet {
		fs := flag.NewFlagSet("up", flag.ExitOnError)
		fs.StringVar(&yggArgs.taskpath, "T", "", "execute template file paths")
		fs.StringVar(&yggArgs.prompt, "P", "", "specify prompt, if not provided by task")
		fs.StringVar(&yggArgs.generator, "G", "openai://gpt-4@localhost:12321", "generator string, {provider}://{model}@{host}:{port}")
		fs.IntVar(&yggArgs.contextWindow, "context-window", 8000, "")
		fs.StringVar(&yggArgs.apiKey, "key", "", "api key by provider models")
		fs.IntVar(&yggArgs.maxIterations, "max-iterations", 0, "max number of automaton to complete task, 0 is the no limit")
		fs.StringVar(&yggArgs.strategy, "S", string(XML), "if a supported format is specified, that format is used")
		fs.BoolVar(&yggArgs.forceFormat, "F", false, "use the fomat specified in serialisation, even if native tools are supported")
		fs.StringVar(&yggArgs.saveTo, "save", "", "at each step, the current system prompts and status data are stored in this file")
		return fs
	})(),
	Exec: exec,
}

func exec(ctx context.Context, args []string) error {
	// setup llm
	options, err := llm.NewLLMOptions(yggArgs.generator, uint32(yggArgs.contextWindow))
	if err != nil {
		return err
	}

	factory, err := llm.NewLLMFactory(options, yggArgs.apiKey)
	if err != nil {
		return err
	}

	// TODO: add embedder for RAG

	// setup task
	tasklet, err := task.GetFromPath(yggArgs.taskpath)
	if err != nil {
		return err
	}

	err = tasklet.Setup(&yggArgs.prompt)
	if err != nil {
		panic(err)
	}

	log.Printf("ygg v%s > 🧬 %s %s", version, yggArgs.generator, tasklet.GetName())

	_, nativeTool := strategyDesicion(StrategyFormat(yggArgs.strategy), yggArgs.forceFormat, factory)
	e := engine.NewEngine(tasklet, factory, uint(yggArgs.maxIterations), nativeTool, yggArgs.saveTo)

	// start
	go e.Start()

	ch := make(chan struct{})
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for {
			select {
			case <-interrupt:
				log.Println("received terminate signal")
				e.Stop()
			case <-e.Done():
				ch <- struct{}{}
				log.Printf("shutdown completed ygg")
			}
		}
	}()
	<-ch

	return nil
}

func strategyDesicion(strategy StrategyFormat, forceFormat bool, factory *llm.LLMFactory) (StrategyFormat, bool) {
	if forceFormat {
		log.Printf("using configured serialization strategy %s\n", strategy)
		return strategy, false
	}
	return strategy, factory.CheckNatvieToolSupport()
}
