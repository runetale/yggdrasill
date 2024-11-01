package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/runetale/notch/engine"
	"github.com/runetale/notch/task"
)

// get from args
var taskpath = "task.yaml"
var prompt = "find the process consuming more ram"

// cli entry point
func main() {
	// setup models

	// setup task
	tasklet, err := task.GetFromPath(taskpath)
	if err != nil {
		panic(err)
	}

	fmt.Printf("notch v0.0.1 🧠 gpt4-o @openai %s", tasklet.Name)

	err = tasklet.Setup(&prompt)
	if err != nil {
		panic(err)
	}

	e := engine.NewEngine()
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
