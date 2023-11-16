package config

import (
	"sync"

	"github.com/caarlos0/env"
)

type Config struct {
	TelegramBotToken string `env:"BOT_TOKEN,required""`
}

var (
	cfg  Config
	once sync.Once
)

func MustLoad(path string) Config {
	once.Do(func() {
		err := env.Parse(&cfg)
		if err != nil {
			panic(err)
		}
	})

	return cfg
}
