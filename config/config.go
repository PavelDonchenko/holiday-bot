package config

import (
	"sync"

	"github.com/caarlos0/env/v10"
)

var (
	cfg  Config
	once sync.Once
)

type (
	Config struct {
		API      API      `env:"API"`
		Telegram Telegram `env:"TELEGRAM"`
	}

	API struct {
		BaseURL        string `env:"BASE_URL"`
		AbstractAPIKey string `env:"ABSTRACT_API_KEY"`
	}

	Telegram struct {
		TelegramBotToken    string `env:"BOT_TOKEN,required"`
		BotDebug            bool   `env:"BOT_DEBUG"`
		UpdateConfigTimeout int    `env:"UPDATE_CONFIG_TIMEOUT"`
	}
)

func MustLoad() Config {
	once.Do(func() {
		if err := env.Parse(&cfg); err != nil {
			panic(err)
		}
	})

	return cfg
}
