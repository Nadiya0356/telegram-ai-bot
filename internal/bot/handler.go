package bot

import (
	"telegram-ai-bot/internal/llm"
	"telegram-ai-bot/internal/session"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler struct {
	llm     llm.LLM
	session *session.Store
}

func New(llm llm.LLM, s *session.Store) *Handler {
	return &Handler{llm: llm, session: s}
}

func (h *Handler) Handle(bot *tgbotapi.BotAPI, update tgbotapi.Update) {

	if update.Message == nil {
		return
	}

	userID := update.Message.From.ID
	text := update.Message.Text

	h.session.Add(userID, "user", text)

	intent, _ := h.llm.Classify(text, h.session.Get(userID))

	var response string

	switch intent {
	case "GREETING":
		response = "Привіт!"
	case "FAQ":
		response = "Це відповідь"
	case "BOOKING":
		response = "Запис зроблено"
	case "FEEDBACK":
		response = "Дякую!"
	default:
		response = "Не зрозуміла"
	}

	h.session.Add(userID, "assistant", response)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
	bot.Send(msg)
}
