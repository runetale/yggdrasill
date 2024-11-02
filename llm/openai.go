// openai api client
package llm

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

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

type OpenAI struct {
	model  string
	client *openai.Client
}

func NewOpenAI(model string, apikey string) LLMClientImpl {
	client := openai.NewClient(apikey)
	return &OpenAI{
		model:  model,
		client: client,
	}
}

func (o *OpenAI) Chat(option *ChatOption) ([]*Invocation, string, error) {
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
	for _, m := range option.GetHistory() {
		switch m.messageType {
		case AGETNT:
			chathistory = append(chathistory, openai.ChatCompletionMessage{
				Role:      openai.ChatMessageRoleAssistant,
				Content:   m.response,
				ToolCalls: nil,
			})
		case FEEDBACK:
			chathistory = append(chathistory, openai.ChatCompletionMessage{
				Role:      openai.ChatMessageRoleUser,
				Content:   m.response,
				ToolCalls: nil,
			})
		}
	}

	// TODO: function tools
	req := openai.ChatCompletionRequest{
		Model:    o.model,
		Messages: chathistory,
		Tools:    nil,
	}

	resp, err := o.client.CreateChatCompletion(
		context.Background(),
		req,
	)

	// TODO: check rate limit
	if err != nil {
		fmt.Println("chat error %v", err)
		return nil, "", err
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

	invocations := make([]*Invocation, 0)
	for _, tool := range toolCalls {
		attributes := make(map[string]string, 0)
		payload := ""

		var result map[string]string
		err := json.Unmarshal([]byte(tool.function.args), &result)
		if err != nil {
			fmt.Println("tool call parsing error %v", err)
			return nil, "", err
		}

		for name, value := range result {
			if name == "payload" {
				payload = value
			} else {
				attributes[name] = value
			}
		}

		in := &Invocation{
			action:     tool.function.name,
			attributes: attributes,
			payload:    payload,
		}

		invocations = append(invocations, in)
	}

	return invocations, content, nil
}
