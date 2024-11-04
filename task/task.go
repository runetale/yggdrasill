// this packege provide user defined task values
package task

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type Task struct {
	name         string         `yaml:"-"`
	folder       string         `yaml:"-"`
	timeout      *time.Duration `yaml:"-"`
	Using        []*string      `yaml:"using"`
	SystemPrompt *string        `yaml:"system_prompt"`
	Prompt       *string        `yaml:"prompt"`
	Guidance     []*string      `yaml:"guidance"`
	Functions    []*Function    `yaml:"functions"`
}

type Function struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Actions     []Action `yaml:"actions"`
}

type Action struct {
	Name           string `yaml:"name"`
	Description    string `yaml:"description"`
	Tool           string `yaml:"tool"`
	MaxShownOutput int    `yaml:"max_shown_output"`
	ExamplePayload string `yaml:"example_payload,omitempty"`
}

func GetFromPath(path string) (*Task, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if info.IsDir() {
		return getFromDir(path)
	}
	return getFromYamlFile(path)
}

func getFromDir(path string) (*Task, error) {
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

func getFromYamlFile(filePath string) (*Task, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("can't read file %v", err)
		return nil, err
	}

	fmt.Println(filePath)
	var tasklet Task
	err = yaml.Unmarshal(data, &tasklet)
	if err != nil {
		log.Fatalf("parsing yaml error %v", err)
		return nil, err
	}

	dir := filepath.Dir(filePath)
	dirName := filepath.Base(dir)
	tasklet.name = dirName
	tasklet.folder = filePath

	return &tasklet, nil
}

// userPromptが無ければ入力を受け付ける
// Promptが設定されていない場合はuserからのpromptを設定する
func (t *Task) Setup(userPrompt *string) error {
	if userPrompt == nil {
		input := t.GetUserInput("enter task > ")
		t.Prompt = &input
		return nil
	}

	if t.Prompt == nil {
		t.Prompt = userPrompt
		return nil
	}

	return errors.New("Setup failed")
}

func (t *Task) GetUsing() []*string {
	return t.Using
}

func (t *Task) GetPrompt() string {
	if t.Prompt != nil {
		return *t.Prompt
	}
	return "no set prompt"
}

func (t *Task) GetSystemPrompt() string {
	if t.SystemPrompt != nil {
		return *t.SystemPrompt
	}
	return "none"
}

func (t *Task) GetMaxHistory() uint {
	return 50
}

func (t *Task) GetGuidance() []*string {
	return t.Guidance
}

// user defined yaml tasks
func (t *Task) GetFunctions() []*Function {
	fs := t.Functions
	log.Println("Called GetFunctions")
	log.Println(fs)
	return t.Functions
}

func (*Task) GetUserInput(prompt string) string {
	log.Print("\n" + prompt)
	log.Print(prompt)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()

	fmt.Println()
	return strings.TrimSpace(input)
}

func (t *Task) ParseVariableExpr(expr string) (string, string, error) {
	if !strings.HasPrefix(expr, "$") {
		return "", "", fmt.Errorf("'%s' is not a valid variable expression", expr)
	}

	varName := strings.TrimPrefix(expr, "$")
	varDefault := ""
	if strings.Contains(varName, "||") {
		parts := strings.SplitN(varName, "||", 2)
		varName = strings.TrimSpace(parts[0])
		varDefault = strings.TrimSpace(parts[1])
	}

	// get from enviroment variables
	if value, exists := os.LookupEnv(varName); exists {
		return varName, value, nil
	}

	// if default value exists
	if varDefault != "" {
		return varName, varDefault, nil
	}

	// user input
	userInput := t.GetUserInput(fmt.Sprintf("\nplease set $%s: ", varName))
	return varName, userInput, nil
}

func (t *Task) GetTimeout() *time.Duration {
	if t.timeout != nil {
		return t.timeout
	}
	return nil
}

func (t *Task) GetName() string {
	return t.name
}
