package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Port int `env:"PORT" envDefault:"8080"`

	TelegramBotToken string `env:"TELEGRAM_BOT_TOKEN,required"`
	TelegramChatID   string `env:"TELEGRAM_CHAT_ID,required"`

	AdminToken string `env:"ADMIN_TOKEN,required"`

	VKConfirmationToken string `env:"VK_CONFIRMATION_TOKEN"`
	VKSecret            string `env:"VK_SECRET"`
}

func Load() (Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return Config{}, fmt.Errorf("parse env: %w", err)
	}

	return cfg, nil
}
