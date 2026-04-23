package llm

type MockClient struct{}

func (m *MockClient) Classify(text string, history []Message) (string, error) {
	if text == "hi" {
		return "GREETING", nil
	}
	return "UNKNOWN", nil
}
