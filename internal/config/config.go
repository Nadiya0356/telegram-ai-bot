package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramToken string
	OpenAIKey     string
}

func Load() Config {
	godotenv.Load()

	return Config{
		TelegramToken: os.Getenv("TELEGRAM_TOKEN"),
		OpenAIKey:     os.Getenv("OPENAI_API_KEY"),
	}
}
