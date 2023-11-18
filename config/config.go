package config

import (
	"sync"

	"github.com/caarlos0/env"
)

type Config struct {
	TelegramBotToken string `env:"BOT_TOKEN,required"`
	AbstractAPIKey   string `env:"ABSTRACT_API_KEY"`
	BaseURL          string `env:"BASE_URL"`
}

var (
	cfg  Config
	once sync.Once
)

func MustLoad() Config {
	once.Do(func() {
		err := env.Parse(&cfg)
		if err != nil {
			panic(err)
		}
	})

	return cfg
}
