package bot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"telegram-ai-bot/internal/llm"
	"telegram-ai-bot/internal/session"
)

type Handler struct {
	llm     llm.LLM
	session *session.Store
}

func New(llm llm.LLM, s *session.Store) *Handler {
	return &Handler{
		llm:     llm,
		session: s,
	}
}

func (h *Handler) Handle(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	userID := update.Message.From.ID
	text := update.Message.Text

	if text == "/start" {
		h.session.Clear(userID)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привіт! Я AI-бот 🤖")
		bot.Send(msg)
		return
	}

	history := h.session.Get(userID)

	intent, err := h.llm.Classify(text, history)
	if err != nil {
		log.Printf("Classify error: %v", err)
		intent = "UNKNOWN"
	}

	log.Printf("userID=%d text=%q intent=%s", userID, text, intent)

	h.session.Add(userID, "user", text)

	var response string

	switch intent {
	case "GREETING":
		response = "Привіт! Чим можу допомогти? 😊"

	case "BOOKING":
		response = "Давайте оформимо запис 📅"

	case "FEEDBACK":
		response = "Дякую за відгук 💬"

	case "FAQ":
		response = "Ось відповідь на ваше питання 📖"

	default:
		resp, err := h.llm.Chat(text, history)
		if err != nil {
			log.Printf("Chat error: %v", err)
			response = "Я поки що не знаю 😅"
		} else if resp == "" {
			log.Println("Chat: empty response")
			response = "Я поки що не знаю 😅"
		} else {
			response = resp
		}
	}

	h.session.Add(userID, "assistant", response)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}
