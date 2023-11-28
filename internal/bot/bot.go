package bot

import (
	"git.foxminded.ua/foxstudent106361/holiday-bot/internal/handler"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"

	"git.foxminded.ua/foxstudent106361/holiday-bot/config"
	"git.foxminded.ua/foxstudent106361/holiday-bot/internal/service"
)

type Bot struct {
	api     *tgbotapi.BotAPI
	service service.Service
	cfg     config.Config
	log     *logrus.Logger
}

func New(api *tgbotapi.BotAPI, cfg config.Config, botService service.Service, log *logrus.Logger) *Bot {
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

	b.createMenu()

	for update := range updates {
		var msg tgbotapi.MessageConfig
		if update.Message != nil {
			if update.Message.Location != nil {
				msg = b.service.UpdateLocation(update.Message)
				msg.ParseMode = tgbotapi.ModeHTML
			} else {
				msg = b.service.UpdateMessage(update.Message)
			}
		}

		if update.CallbackQuery != nil {
			msg = b.service.UpdateCallback(update.CallbackQuery)
		}

		_, err := b.api.Send(msg)
		if err != nil {
			b.log.Errorf("error sending message, err: %v", err)
			return
		}
	}
}

func (b *Bot) createMenu() {
	cfg := tgbotapi.NewSetMyCommands(
		tgbotapi.BotCommand{
			Command:     handler.StartMenu,
			Description: "Show menu",
		})

	_, _ = b.api.Request(cfg)
}
