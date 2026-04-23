package session

import "testing"

func TestSession(t *testing.T) {
	s := New()
	s.Add(1, "user", "hi")

	if len(s.Get(1)) != 1 {
		t.Error("Session not working")
	}
}
