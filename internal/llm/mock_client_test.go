package llm

import (
	"strings"
	"testing"
)

func TestMockClient_Classify_Greeting(t *testing.T) {
	m := &MockClient{}
	intent, err := m.Classify("hi", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if intent != "GREETING" {
		t.Errorf("expected GREETING, got %s", intent)
	}
}

func TestMockClient_Classify_Unknown(t *testing.T) {
	m := &MockClient{}
	intent, err := m.Classify("what is the meaning of life?", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if intent != "UNKNOWN" {
		t.Errorf("expected UNKNOWN, got %s", intent)
	}
}

func TestMockClient_Chat_ReturnsResponse(t *testing.T) {
	m := &MockClient{}
	resp, err := m.Chat("tell me something", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == "" {
		t.Error("expected non-empty response")
	}
}

func TestMockClient_Chat_ContainsInput(t *testing.T) {
	m := &MockClient{}
	input := "hello world"
	resp, err := m.Chat(input, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(resp, input) {
		t.Errorf("expected response to contain input %q, got %q", input, resp)
	}
}

func TestMockClient_Classify_WithHistory(t *testing.T) {
	m := &MockClient{}
	history := []Message{
		{Role: "user", Content: "previous message"},
		{Role: "assistant", Content: "previous response"},
	}
	intent, err := m.Classify("hi", history)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if intent == "" {
		t.Error("expected non-empty intent")
	}
}

// Table-driven тест для різних варіантів вхідного тексту
func TestMockClient_Classify_TableDriven(t *testing.T) {
	m := &MockClient{}

	tests := []struct {
		input    string
		expected string
	}{
		{"hi", "GREETING"},
		{"hello", "UNKNOWN"}, // mock реагує тільки на "hi"
		{"book a table", "UNKNOWN"},
		{"", "UNKNOWN"},
		{"hi there how are you", "UNKNOWN"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := m.Classify(tt.input, nil)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.expected {
				t.Errorf("Classify(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}
