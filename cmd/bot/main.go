package main

import (
	"log"

	"telegram-ai-bot/internal/bot"
	"telegram-ai-bot/internal/config"
	"telegram-ai-bot/internal/llm"
	"telegram-ai-bot/internal/session"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	cfg := config.Load()

	log.Println("Starting AI bot...")

	b, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Fatal(err)
	}

	llmClient := llm.New(cfg.OpenRouterKey)
	store := session.New()
	handler := bot.New(llmClient, store)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	for update := range b.GetUpdatesChan(u) {
		handler.Handle(b, update)
	}
}
