package bot

import (
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"telegram-ai-bot/internal/llm"
	"telegram-ai-bot/internal/session"
)

// makeUpdate створює тестовий update з повідомленням
func makeUpdate(userID int64, chatID int64, text string) tgbotapi.Update {
	return tgbotapi.Update{
		Message: &tgbotapi.Message{
			From: &tgbotapi.User{ID: userID},
			Chat: &tgbotapi.Chat{ID: chatID},
			Text: text,
		},
	}
}

func TestHandler_IgnoresEmptyUpdate(t *testing.T) {
	h := New(&llm.MockClient{}, session.New())
	// nil Message — не має панікувати
	h.Handle(nil, tgbotapi.Update{Message: nil})
}

func TestHandler_Start_ClearsSession(t *testing.T) {
	store := session.New()
	store.Add(1, "user", "old message")

	h := New(&llm.MockClient{}, store)

	update := makeUpdate(1, 100, "/start")

	// bot.Send паніккує з nil — перехоплюємо, нас цікавить тільки стан сесії
	func() {
		defer func() { recover() }()
		h.Handle(nil, update)
	}()

	if len(store.Get(1)) != 0 {
		t.Error("session should be cleared after /start")
	}
}

func TestHandler_Greeting_AddsToSession(t *testing.T) {
	store := session.New()
	h := New(&llm.MockClient{}, store)

	func() {
		defer func() { recover() }()
		h.Handle(nil, makeUpdate(1, 100, "hi"))
	}()

	history := store.Get(1)
	if len(history) == 0 {
		t.Fatal("expected messages in session after handling")
	}
	if history[0].Role != "user" || history[0].Content != "hi" {
		t.Errorf("first message should be user 'hi', got %+v", history[0])
	}
}

func TestHandler_SessionIsolatedPerUser(t *testing.T) {
	store := session.New()
	h := New(&llm.MockClient{}, store)

	func() {
		defer func() { recover() }()
		h.Handle(nil, makeUpdate(1, 100, "hi"))
	}()
	func() {
		defer func() { recover() }()
		h.Handle(nil, makeUpdate(2, 200, "hi"))
	}()

	if len(store.Get(1)) == 0 {
		t.Error("user 1 should have session data")
	}
	if len(store.Get(2)) == 0 {
		t.Error("user 2 should have session data")
	}
}

func TestHandler_AssistantResponseSavedToSession(t *testing.T) {
	store := session.New()
	h := New(&llm.MockClient{}, store)

	func() {
		defer func() { recover() }()
		h.Handle(nil, makeUpdate(1, 100, "hi"))
	}()

	history := store.Get(1)
	if len(history) < 2 {
		t.Fatalf("expected user + assistant messages, got %d", len(history))
	}
	if history[1].Role != "assistant" {
		t.Errorf("second message should be assistant, got %q", history[1].Role)
	}
}
