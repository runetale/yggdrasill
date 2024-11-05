package notch

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/peterbourgon/ff/v2/ffcli"
	"github.com/runetale/notch/engine"
	"github.com/runetale/notch/llm"
	"github.com/runetale/notch/task"
)

const version = "0.0.1"

var notchArgs struct {
	taskpath      string
	prompt        string
	generator     string
	contextWindow int
	apiKey        string
	maxIterations int
	strategy      string
	forceFormat   bool
}

type StrategyFormat string

const (
	XML StrategyFormat = "xml"
)

var NotchCmd = &ffcli.Command{
	Name:       "up",
	ShortUsage: "up [flags]",
	ShortHelp:  "command to start notch",
	FlagSet: (func() *flag.FlagSet {
		fs := flag.NewFlagSet("up", flag.ExitOnError)
		fs.StringVar(&notchArgs.taskpath, "T", "", "execute template file paths")
		fs.StringVar(&notchArgs.prompt, "P", "", "specify prompt, if not provided by task")
		fs.StringVar(&notchArgs.generator, "G", "openai://gpt-4@localhost:12321", "generator string, {provider}://{model}@{host}:{port}")
		fs.IntVar(&notchArgs.contextWindow, "context-window", 8000, "")
		fs.StringVar(&notchArgs.apiKey, "key", "", "api key by provider models")
		fs.IntVar(&notchArgs.maxIterations, "max-iterations", 0, "max number of automaton to complete task, 0 is the no limit")
		fs.StringVar(&notchArgs.strategy, "S", string(XML), "if a supported format is specified, that format is used")
		fs.BoolVar(&notchArgs.forceFormat, "F", false, "use the fomat specified in serialisation, even if native tools are supported")
		return fs
	})(),
	Exec: exec,
}

func exec(ctx context.Context, args []string) error {
	// setup llm
	options, err := llm.NewLLMOptions(notchArgs.generator, uint32(notchArgs.contextWindow))
	if err != nil {
		return err
	}

	factory, err := llm.NewLLMFactory(options, notchArgs.apiKey)
	if err != nil {
		return err
	}

	// TODO: add embedder for RAG

	// setup task
	tasklet, err := task.GetFromPath(notchArgs.taskpath)
	if err != nil {
		return err
	}

	err = tasklet.Setup(&notchArgs.prompt)
	if err != nil {
		panic(err)
	}

	log.Printf("notch v%s > 🧬 %s %s", version, notchArgs.generator, tasklet.GetName())

	_, nativeTool := strategyDesicion(StrategyFormat(notchArgs.strategy), notchArgs.forceFormat, factory)
	e := engine.NewEngine(tasklet, factory, uint(notchArgs.maxIterations), nativeTool)

	// start
	go e.Start()

	interrupt := make(chan os.Signal, 1)
	ch := make(chan struct{})
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for {
			select {
			case <-interrupt:
				close(ch)
				return
			case <-e.Done():
				close(ch)
				return
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
