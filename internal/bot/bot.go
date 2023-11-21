package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"git.foxminded.ua/foxstudent106361/holiday-bot/config"
	"git.foxminded.ua/foxstudent106361/holiday-bot/internal/service"
	"git.foxminded.ua/foxstudent106361/holiday-bot/pkg/logging"
)

type Bot struct {
	api     *tgbotapi.BotAPI
	service service.Service
	cfg     config.Config
	log     logging.Logger
}

func New(api *tgbotapi.BotAPI, cfg config.Config, botService service.Service, log logging.Logger) *Bot {
	return &Bot{
		api:     api,
		cfg:     cfg,
		log:     log,
		service: botService,
	}
}

func (b *Bot) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = b.cfg.Telegram.UpdateConfigTimeout

	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			msg := b.service.HandleMessage(update.Message)

			_, err := b.api.Send(msg)
			if err != nil {
				b.log.Errorf("error sending message, err: %v", err)
				return
			}
		}
	}
}
