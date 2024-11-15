// openai api client
package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/runetale/yggdrasill/engine/chat"
	"github.com/runetale/yggdrasill/engine/namespace"
	"github.com/sashabaranov/go-openai"
)

type ToolFunctionParameterProperty struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

type ToolFunctionParameter struct {
	Type       string                                   `json:"type"`
	Required   []string                                 `json:"required"`
	Properties map[string]ToolFunctionParameterProperty `json:"properties"`
}

const (
	O1Mini                = "o1-mini"
	O1Mini20240912        = "o1-mini-2024-09-12"
	O1Preview             = "o1-preview"
	O1Preview20240912     = "o1-preview-2024-09-12"
	GPT432K0613           = "gpt-4-32k-0613"
	GPT432K0314           = "gpt-4-32k-0314"
	GPT432K               = "gpt-4-32k"
	GPT40613              = "gpt-4-0613"
	GPT40314              = "gpt-4-0314"
	GPT4o                 = "gpt-4o"
	GPT4o20240513         = "gpt-4o-2024-05-13"
	GPT4o20240806         = "gpt-4o-2024-08-06"
	GPT4oLatest           = "chatgpt-4o-latest"
	GPT4oMini             = "gpt-4o-mini"
	GPT4oMini20240718     = "gpt-4o-mini-2024-07-18"
	GPT4Turbo             = "gpt-4-turbo"
	GPT4Turbo20240409     = "gpt-4-turbo-2024-04-09"
	GPT4Turbo0125         = "gpt-4-0125-preview"
	GPT4Turbo1106         = "gpt-4-1106-preview"
	GPT4TurboPreview      = "gpt-4-turbo-preview"
	GPT4VisionPreview     = "gpt-4-vision-preview"
	GPT4                  = "gpt-4"
	GPT3Dot5Turbo0125     = "gpt-3.5-turbo-0125"
	GPT3Dot5Turbo1106     = "gpt-3.5-turbo-1106"
	GPT3Dot5Turbo0613     = "gpt-3.5-turbo-0613"
	GPT3Dot5Turbo0301     = "gpt-3.5-turbo-0301"
	GPT3Dot5Turbo16K      = "gpt-3.5-turbo-16k"
	GPT3Dot5Turbo16K0613  = "gpt-3.5-turbo-16k-0613"
	GPT3Dot5Turbo         = "gpt-3.5-turbo"
	GPT3Dot5TurboInstruct = "gpt-3.5-turbo-instruct"
)

type ToolCall struct {
	id       string
	function Function
	theType  string
}

type Function struct {
	name string
	args string
}

type OpenAIClient struct {
	model  string
	client *openai.Client
	url    string
	port   uint16
}

func NewOpenAIClient(model string, apikey string, url string, port uint16) LLMClientImpl {
	client := openai.NewClient(apikey)
	return &OpenAIClient{
		model:  model,
		client: client,
		url:    url,
		port:   port,
	}
}

func (o *OpenAIClient) Chat(option *chat.ChatOption, nativeSupport bool, namespaces []*namespace.Namespace) ([]*chat.Invocation, string) {
	chathistory := []openai.ChatCompletionMessage{
		{
			Role:      openai.ChatMessageRoleSystem,
			Content:   option.GetSystemPrompt(),
			ToolCalls: nil,
		},
		{
			Role:      openai.ChatMessageRoleUser,
			Content:   option.GetPrompt(),
			ToolCalls: nil,
		},
	}

	// add chat history
	if option.GetHistory() != nil {
		for _, m := range option.GetHistory() {
			switch m.MessageType {
			case chat.AGETNT:
				chathistory = append(chathistory, openai.ChatCompletionMessage{
					Role:      openai.ChatMessageRoleAssistant,
					Content:   *m.Response,
					ToolCalls: nil,
				})
			case chat.FEEDBACK:
				chathistory = append(chathistory, openai.ChatCompletionMessage{
					Role:      openai.ChatMessageRoleUser,
					Content:   *m.Response,
					ToolCalls: nil,
				})
			}
		}
	}

	// add native tools function
	// setting each namespaces actions,
	// todo: toolsの内容が正しいか？
	tools := []openai.Tool{}
	if nativeSupport {
		for _, group := range namespaces {
			for _, action := range group.GetActions() {
				required := []string{}
				properties := map[string]ToolFunctionParameterProperty{}

				if action.ExamplePayload() != nil {
					required = append(required, "payload")
					properties["payload"] = ToolFunctionParameterProperty{
						Type:        "string",
						Description: "Main function argument.",
					}
				}

				for key, _ := range action.ExampleAttributes() {
					required = append(required, key)
					properties[key] = ToolFunctionParameterProperty{
						Type:        "string",
						Description: key,
					}
				}

				function := &openai.FunctionDefinition{
					Name:        action.Name(),
					Description: action.Description(),
					Parameters: ToolFunctionParameter{
						Type:       "object",
						Required:   required,
						Properties: properties,
					},
				}

				tools = append(tools, openai.Tool{
					Type:     "function",
					Function: function,
				})
			}
		}
	}

	// request to chat
	req := openai.ChatCompletionRequest{
		Model:    o.model,
		Messages: chathistory,
		Tools:    tools,
	}

	resp, err := o.client.CreateChatCompletion(
		context.Background(),
		req,
	)

	// TODO: check rate limit, retry to chat
	if err != nil {
		fmt.Printf("chat error %v\n", err)
		panic(err)
	}

	// add invocations
	choice := resp.Choices[0]
	content := choice.Message.Content
	toolCalls := make([]ToolCall, 0)
	for _, tool := range choice.Message.ToolCalls {
		toolCalls = append(toolCalls, ToolCall{
			id: tool.ID,
			function: Function{
				name: tool.Function.Name,
				args: tool.Function.Arguments,
			},
			theType: string(tool.Type),
		})
	}

	invocations := make([]*chat.Invocation, 0)
	for _, tool := range toolCalls {
		attributes := make(map[string]string, 0)
		payload := ""

		var result map[string]string
		err := json.Unmarshal([]byte(tool.function.args), &result)
		if err != nil {
			fmt.Printf("tool call parsing error %v\n", err)
		}

		for name, value := range result {
			if name == "payload" {
				payload = value
			} else {
				attributes[name] = value
			}
		}

		in := chat.NewInvocation(tool.function.name, attributes, &payload)

		invocations = append(invocations, in)
	}

	return invocations, content
}

func (o *OpenAIClient) CheckNatvieToolSupport() bool {
	chathistory := []openai.ChatCompletionMessage{
		{
			Role:      openai.ChatMessageRoleSystem,
			Content:   "You are an helpful assistant.",
			ToolCalls: nil,
		},
		{
			Role:      openai.ChatMessageRoleUser,
			Content:   "Call the test function.",
			ToolCalls: nil,
		},
	}

	functionTools := []openai.Tool{}
	functionTools = append(functionTools, openai.Tool{
		Type: "function",
		Function: &openai.FunctionDefinition{
			Name:        "test",
			Description: "This is a test function.",
			Parameters:  nil,
		},
	})

	req := openai.ChatCompletionRequest{
		Model:    o.model,
		Messages: chathistory,
		Tools:    functionTools,
	}

	resp, err := o.client.CreateChatCompletion(
		context.Background(),
		req,
	)

	if err != nil {
		log.Printf("error check native tool support request error %s", err.Error())
		return false
	}

	if resp.Choices != nil {
		m := resp.Choices[0].Message
		if m.ToolCalls != nil {
			log.Printf("using native tools by %s", o.model)
			return true
		}
	}

	log.Println("using original yggdrasill system prompt")
	return false
}
