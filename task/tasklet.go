package task

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type Tasklet struct {
	Name         string
	Folder       string
	Using        []string `yaml:"using"`
	SystemPrompt string   `yaml:"system_prompt"`
	Prompt       *string  `yaml:"prompt"`
	Guidance     []string `yaml:"guidance"`
}

func GetFromPath(path string) (*Tasklet, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if info.IsDir() {
		return getFromDir(path)
	}
	return getFromYamlFile(path)
}

func getFromDir(path string) (*Tasklet, error) {
	filePath := filepath.Join(path, "task.yaml")
	_, err := os.Stat(filePath)
	if err == nil {
		return getFromYamlFile(filePath)
	}
	if os.IsNotExist(err) {
		return nil, err
	}

	return nil, err
}

func getFromYamlFile(filePath string) (*Tasklet, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("can't read file %v", err)
		return nil, err
	}

	var tasklet Tasklet
	err = yaml.Unmarshal(data, &tasklet)
	if err != nil {
		log.Fatalf("parsing yaml error %v", err)
		return nil, err
	}

	dir := filepath.Dir(filePath)
	dirName := filepath.Base(dir)
	tasklet.Name = dirName
	tasklet.Folder = filePath

	return &tasklet, nil

}

// userPromptが無ければ入力を受け付ける
// Promptが設定されていない場合はuserからのpromptを設定する
func (t *Tasklet) Setup(userPrompt *string) error {
	if userPrompt == nil {
		input := getUserInput("enter task > ")
		t.Prompt = &input
		return nil
	}

	if t.Prompt == nil {
		t.Prompt = userPrompt
		return nil
	}

	return errors.New("Setup failed")
}
func getUserInput(prompt string) string {
	fmt.Print("\n" + prompt)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()

	fmt.Println()
	return strings.TrimSpace(input)
}
