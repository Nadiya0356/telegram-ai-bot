package session

import (
	"testing"
)

func TestAdd_SingleMessage(t *testing.T) {
	s := New()
	s.Add(1, "user", "hi")

	msgs := s.Get(1)
	if len(msgs) != 1 {
		t.Fatalf("expected 1 message, got %d", len(msgs))
	}
	if msgs[0].Role != "user" || msgs[0].Content != "hi" {
		t.Errorf("unexpected message: %+v", msgs[0])
	}
}

func TestAdd_MultipleMessages(t *testing.T) {
	s := New()
	s.Add(1, "user", "hello")
	s.Add(1, "assistant", "hi there")
	s.Add(1, "user", "how are you?")

	msgs := s.Get(1)
	if len(msgs) != 3 {
		t.Fatalf("expected 3 messages, got %d", len(msgs))
	}
}

func TestAdd_LimitTo10Messages(t *testing.T) {
	s := New()
	for i := 0; i < 15; i++ {
		s.Add(1, "user", "message")
	}

	msgs := s.Get(1)
	if len(msgs) != 10 {
		t.Errorf("expected 10 messages (limit), got %d", len(msgs))
	}
}

func TestGet_EmptySession(t *testing.T) {
	s := New()
	msgs := s.Get(999)
	if msgs != nil && len(msgs) != 0 {
		t.Errorf("expected empty slice for unknown user, got %v", msgs)
	}
}

func TestGet_IsolatedPerUser(t *testing.T) {
	s := New()
	s.Add(1, "user", "hello from user 1")
	s.Add(2, "user", "hello from user 2")

	u1 := s.Get(1)
	u2 := s.Get(2)

	if len(u1) != 1 || u1[0].Content != "hello from user 1" {
		t.Errorf("user 1 session corrupted: %+v", u1)
	}
	if len(u2) != 1 || u2[0].Content != "hello from user 2" {
		t.Errorf("user 2 session corrupted: %+v", u2)
	}
}

func TestClear(t *testing.T) {
	s := New()
	s.Add(1, "user", "hi")
	s.Add(1, "assistant", "hello")
	s.Clear(1)

	msgs := s.Get(1)
	if len(msgs) != 0 {
		t.Errorf("expected empty after Clear, got %d messages", len(msgs))
	}
}

func TestClear_OnlyAffectsOneUser(t *testing.T) {
	s := New()
	s.Add(1, "user", "hi")
	s.Add(2, "user", "hello")
	s.Clear(1)

	if len(s.Get(1)) != 0 {
		t.Error("user 1 should be cleared")
	}
	if len(s.Get(2)) != 1 {
		t.Error("user 2 should not be affected by Clear of user 1")
	}
}

func TestAdd_PreservesOrder(t *testing.T) {
	s := New()
	s.Add(1, "user", "first")
	s.Add(1, "assistant", "second")
	s.Add(1, "user", "third")

	msgs := s.Get(1)
	expected := []string{"first", "second", "third"}
	for i, exp := range expected {
		if msgs[i].Content != exp {
			t.Errorf("position %d: expected %q, got %q", i, exp, msgs[i].Content)
		}
	}
}

func TestAdd_LimitKeepsLatestMessages(t *testing.T) {
	s := New()
	for i := 0; i < 12; i++ {
		s.Add(1, "user", string(rune('a'+i))) // a, b, c, ... l
	}

	msgs := s.Get(1)
	// має залишити останні 10: c, d, e, f, g, h, i, j, k, l
	if msgs[0].Content != string(rune('a'+2)) {
		t.Errorf("expected oldest kept message to be 'c', got %q", msgs[0].Content)
	}
	if msgs[9].Content != string(rune('a'+11)) {
		t.Errorf("expected newest message to be 'l', got %q", msgs[9].Content)
	}
}
