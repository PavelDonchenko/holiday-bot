package bot

import (
	"git.foxminded.ua/foxstudent106361/holiday-bot/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const updateConfigTimeout = 60

type Bot struct {
	api     *tgbotapi.BotAPI
	fetcher Fetcher
	cfg     config.Config
}

func New(api *tgbotapi.BotAPI, cfg config.Config, fetcher Fetcher) *Bot {
	return &Bot{
		api:     api,
		cfg:     cfg,
		fetcher: fetcher,
	}
}

func (b *Bot) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = updateConfigTimeout

	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			b.handleMessage(update.Message)
		}
	}
}

func (b *Bot) handleMessage(message *tgbotapi.Message) {
	switch message.Text {
	case "🇺🇸 USA", "🇬🇧 UK", "🇨🇦 Canada", "🇫🇷 France", "🇩🇪 Germany", "🇯🇵 Japan":
		b.handleGetHolidays(message)
	default:
		b.handleFlags(message)
	}
}
