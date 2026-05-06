package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Message represents a single turn in a conversation.
type Message struct {
	Role    string
	Content string
}

// LLM describes the capabilities of a language model client.
type LLM interface {
	Classify(text string, history []Message) (string, error)
	Chat(text string, history []Message) (string, error)
}

// ORClient is an OpenRouter-backed implementation of LLM.
type ORClient struct {
	apiKey     string
	httpClient *http.Client
	model      string
	baseURL    string
}

// New creates a new ORClient with sensible defaults.
func New(apiKey string) *ORClient {
	return &ORClient{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
		model:   "llama-3.1-8b-instant",
		baseURL: "https://api.groq.com/openai/v1/chat/completions",
	}
}

// orMessage is the wire format expected by the OpenRouter API.
type orMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// orRequest is the full request body sent to the API.
type orRequest struct {
	Model    string      `json:"model"`
	Messages []orMessage `json:"messages"`
}

// orResponse is the partial shape of the API response we care about.
type orResponse struct {
	Error *struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
	} `json:"error,omitempty"`
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// complete sends messages to the OpenRouter API and returns the assistant reply.
func (c *ORClient) complete(messages []orMessage) (string, error) {
	payload := orRequest{
		Model:    c.model,
		Messages: messages,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, c.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API returned status %d: %s", resp.StatusCode, body)
	}

	var result orResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("unmarshal response: %w", err)
	}

	if result.Error != nil {
		return "", fmt.Errorf("API error (code %d): %s", result.Error.Code, result.Error.Message)
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("API returned no choices")
	}

	return result.Choices[0].Message.Content, nil
}

// convertHistory maps domain Messages to the wire format.
func convertHistory(history []Message) []orMessage {
	out := make([]orMessage, len(history))
	for i, m := range history {
		out[i] = orMessage{Role: m.Role, Content: m.Content}
	}
	return out
}

// Classify uses the AI to determine the intent of `text`.
// It returns a short label such as "GREETING", "QUESTION", "COMPLAINT", etc.
func (c *ORClient) Classify(text string, history []Message) (string, error) {
	systemPrompt := orMessage{
		Role: "system",
		Content: `You are an intent classifier. 
Given a user message, respond with ONLY one of these labels — nothing else:
GREETING, QUESTION, COMPLAINT, REQUEST, FAREWELL, UNKNOWN.`,
	}

	messages := []orMessage{systemPrompt}
	messages = append(messages, convertHistory(history)...)
	messages = append(messages, orMessage{Role: "user", Content: text})

	label, err := c.complete(messages)
	if err != nil {
		return "", fmt.Errorf("classify: %w", err)
	}

	return label, nil
}

// Chat sends a user message together with the full conversation history
// and returns the assistant's reply.
func (c *ORClient) Chat(text string, history []Message) (string, error) {
	messages := convertHistory(history)
	messages = append(messages, orMessage{Role: "user", Content: text})

	reply, err := c.complete(messages)
	if err != nil {
		return "", fmt.Errorf("chat: %w", err)
	}

	return reply, nil
}
