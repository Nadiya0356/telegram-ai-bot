package llm

import (
	"context"
	"encoding/json"

	openai "github.com/sashabaranov/go-openai"
)

type Message struct {
	Role    string
	Content string
}

type LLM interface {
	Classify(text string, history []Message) (string, error)
}

type OpenAIClient struct {
	api *openai.Client
}

func New(apiKey string) *OpenAIClient {
	return &OpenAIClient{
		api: openai.NewClient(apiKey),
	}
}

func (c *OpenAIClient) Classify(text string, history []Message) (string, error) {

	messages := []openai.ChatCompletionMessage{
		{
			Role:    "system",
			Content: "Return ONLY JSON: {\"intent\":\"GREETING|FAQ|BOOKING|FEEDBACK|UNKNOWN\"}",
		},
	}

	for _, h := range history {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    h.Role,
			Content: h.Content,
		})
	}

	messages = append(messages, openai.ChatCompletionMessage{
		Role:    "user",
		Content: text,
	})

	resp, err := c.api.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model:    "gpt-4o-mini",
		Messages: messages,
	})

	if err != nil {
		return "UNKNOWN", nil
	}

	var result map[string]string
	err = json.Unmarshal([]byte(resp.Choices[0].Message.Content), &result)

	if err != nil {
		return "UNKNOWN", nil
	}

	return result["intent"], nil
}
