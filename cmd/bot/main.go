package main

import (
	"telegram-ai-bot/internal/bot"
	"telegram-ai-bot/internal/config"
	"telegram-ai-bot/internal/llm"
	"telegram-ai-bot/internal/session"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	cfg := config.Load()

	b, _ := tgbotapi.NewBotAPI(cfg.TelegramToken)

	llmClient := llm.New(cfg.OpenAIKey)
	store := session.New()
	handler := bot.New(llmClient, store)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.GetUpdatesChan(u)

	for update := range updates {
		handler.Handle(b, update)
	}
}
