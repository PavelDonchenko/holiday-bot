package config

import (
	"sync"

	"github.com/caarlos0/env/v10"
)

var (
	cfg  Config
	once sync.Once
)

type Config struct {
	API      API      `env:"API"`
	Telegram Telegram `env:"TELEGRAM"`
}

type API struct {
	BaseAbstractURL string `env:"BASE_ABSTRACT_URL"`
	BaseWeatherURL  string `env:"BASE_WEATHER_URL"`
	AbstractAPIKey  string `env:"ABSTRACT_API_KEY"`
	WeatherAPIKey   string `env:"WEATHER_API_KEY"`
}

type Telegram struct {
	TelegramBotToken    string `env:"BOT_TOKEN,required"`
	BotDebug            bool   `env:"BOT_DEBUG"`
	UpdateConfigTimeout int    `env:"UPDATE_CONFIG_TIMEOUT"`
}

func MustLoad() Config {
	once.Do(func() {
		if err := env.Parse(&cfg); err != nil {
			panic(err)
		}
	})

	return cfg
}
