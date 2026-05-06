package llm

type MockClient struct{}

func (m *MockClient) Classify(text string, history []Message) (string, error) {
	if text == "hi" {
		return "GREETING", nil
	}
	return "UNKNOWN", nil
}

// Chat потрібен щоб MockClient реалізував інтерфейс LLM
func (m *MockClient) Chat(text string, history []Message) (string, error) {
	return "Mock відповідь: " + text, nil
}
