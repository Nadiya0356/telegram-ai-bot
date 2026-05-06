package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramToken string
	OpenRouterKey string
}

func Load() Config {
	godotenv.Load()

	return Config{
		TelegramToken: os.Getenv("TELEGRAM_TOKEN"),
		OpenRouterKey: os.Getenv("OPENROUTER_API_KEY"),
	}
}
