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
	HolidayAPI  HolidayAPI
	WeatherAPI  WeatherAPI
	Telegram    Telegram
	Mongo       Mongo
	Worker      Worker
	Application Aplication
}

type HolidayAPI struct {
	BaseHolidayURL string `env:"BASE_HOLIDAY_URL"`
	HolidayAPIKey  string `env:"HOLIDAY_API_KEY"`
}

type WeatherAPI struct {
	BaseWeatherURL string `env:"BASE_WEATHER_URL"`
	WeatherAPIKey  string `env:"WEATHER_API_KEY"`
}

type Telegram struct {
	TelegramBotToken    string `env:"BOT_TOKEN,required"`
	BotDebug            bool   `env:"BOT_DEBUG"`
	UpdateConfigTimeout int    `env:"UPDATE_CONFIG_TIMEOUT"`
}

type Mongo struct {
	URL        string `env:"MONGODB_LOCAL_URI"`
	Database   string `env:"DATABASE"`
	Collection string `env:"COLLECTION"`
}

type Worker struct {
	WorkerDuration int `env:"WORKER_DURATION"`
}

type Aplication struct {
	Run string `env:"RUN"`
}

func MustLoad() Config {
	once.Do(func() {
		if err := env.Parse(&cfg); err != nil {
			panic(err)
		}
	})

	return cfg
}
