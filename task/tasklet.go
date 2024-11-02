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
	Using        []*string  `yaml:"using"`
	SystemPrompt string     `yaml:"system_prompt"`
	Prompt       *string    `yaml:"prompt"`
	Guidance     []string   `yaml:"guidance"`
	Functions    []Function `yaml:"functions"`
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

func (t *Tasklet) GetUsing() []*string {
	return t.Using
}

func (t *Tasklet) GetPrompt() *string {
	return t.Prompt
}

func (t *Tasklet) GetSystemPrompt() string {
	return t.SystemPrompt
}

// user defined yaml tasks
func (t *Tasklet) GetFunctions() []Function {
	fs := t.Functions
	fmt.Println("Called GetFunctions")
	fmt.Println(fs)
	return t.Functions
}

func (*Tasklet) GetUserInput(prompt string) string {
	fmt.Print("\n" + prompt)
	fmt.Print(prompt) // プロンプトを表示

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()

	fmt.Println()
	return strings.TrimSpace(input)
}

func (t *Tasklet) ParseVariableExpr(expr string) (string, string, error) {
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
