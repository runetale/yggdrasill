package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/runetale/notch/engine"
	"github.com/runetale/notch/llm"
	"github.com/runetale/notch/task"
)

// TODO: fix get from args
var taskpath = "task.yaml"
var prompt = "find the process consuming more ram"
var generator = "openai://gpt-4@localhost:12321"
var contextWindow uint32 = 8000
var apiKey = ""

func main() {
	// setup llm
	options, err := llm.NewLLMOptions(generator, contextWindow)
	if err != nil {
		panic(err)
	}

	llm, err := llm.NewLLMClient(options, apiKey)
	if err != nil {
		panic(err)
	}

	// TODO: add embedder for RAG

	// setup task
	tasklet, err := task.GetFromPath(taskpath)
	if err != nil {
		panic(err)
	}

	err = tasklet.Setup(&prompt)
	if err != nil {
		panic(err)
	}

	fmt.Printf("notch v0.0.1 🧠 gpt4-o @openai %s", tasklet.Name)

	e := engine.NewEngine(tasklet, llm)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for {
			select {
			case <-interrupt:
				return
			case <-e.Done():
				return
			}
		}
	}()

	// start
	go e.Start()
}
