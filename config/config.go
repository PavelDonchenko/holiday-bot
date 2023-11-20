package config

import (
	"sync"

	"github.com/caarlos0/env"
)

var (
	cfg  Config
	once sync.Once
)

type Config struct {
	Telegram struct {
		TelegramBotToken    string `env:"BOT_TOKEN,required"`
		BotDebug            bool   `env:"BOT_DEBUG"`
		UpdateConfigTimeout int    `env:"UPDATE_CONFIG_TIMEOUT"`
	}
	API struct {
		BaseURL        string `env:"BASE_URL"`
		AbstractAPIKey string `env:"ABSTRACT_API_KEY"`
	}
}

func MustLoad() Config {
	once.Do(func() {
		err := env.Parse(&cfg)
		if err != nil {
			panic(err)
		}
	})

	return cfg
}
