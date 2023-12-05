package bot

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"

	"git.foxminded.ua/foxstudent106361/holiday-bot/internal/handler"
	"git.foxminded.ua/foxstudent106361/holiday-bot/internal/model"

	"git.foxminded.ua/foxstudent106361/holiday-bot/config"
	"git.foxminded.ua/foxstudent106361/holiday-bot/internal/service"
)

type Bot struct {
	api                 *tgbotapi.BotAPI
	service             service.Service
	cfg                 config.Config
	log                 *logrus.Logger
	userState           model.State
	subscribeCommands   chan tgbotapi.Update
	updateTimeCommands  chan tgbotapi.Update
	regularCommands     chan tgbotapi.Update
	holidayCommands     chan tgbotapi.Update
	unsubscribeCommands chan tgbotapi.Update
	weatherCommands     chan tgbotapi.Update
}

func New(api *tgbotapi.BotAPI, cfg config.Config, botService service.Service, log *logrus.Logger) *Bot {
	return &Bot{
		api:                 api,
		cfg:                 cfg,
		log:                 log,
		service:             botService,
		userState:           make(map[int64]model.ActiveFlags),
		regularCommands:     make(chan tgbotapi.Update),
		holidayCommands:     make(chan tgbotapi.Update),
		subscribeCommands:   make(chan tgbotapi.Update),
		unsubscribeCommands: make(chan tgbotapi.Update),
		updateTimeCommands:  make(chan tgbotapi.Update),
		weatherCommands:     make(chan tgbotapi.Update),
	}
}

func (b *Bot) Run(ctx context.Context) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = b.cfg.Telegram.UpdateConfigTimeout

	updates := b.api.GetUpdatesChan(u)

	b.createMenu()

	var channelToSend chan tgbotapi.Update
	msgChan := make(chan tgbotapi.MessageConfig)

	go func() {
		for update := range updates {
			chatID := getChatID(update)

			// need to fill flags with false value from the start and when user change flow in the middle of processing
			if update.Message != nil {
				if update.Message.IsCommand() {
					b.userState[chatID] = model.ActiveFlags{}
				}
			}

			if b.userState[chatID].UnsubscribeActiveFlag {
				channelToSend = b.unsubscribeCommands
			} else if b.userState[chatID].SubscribeActiveFlag {
				channelToSend = b.subscribeCommands
			} else if b.userState[chatID].HolidayActiveFlag {
				channelToSend = b.holidayCommands
			} else if b.userState[chatID].WeatherActiveFlag {
				channelToSend = b.weatherCommands
			} else if b.userState[chatID].UpdateTimeActiveFlag {
				channelToSend = b.updateTimeCommands
			} else {
				channelToSend = b.regularCommands
			}

			channelToSend <- update
		}
	}()

	go func() {
		for update := range b.regularCommands {
			b.service.UpdateRegularCommand(update.Message, update.Message.Chat.ID, b.userState, msgChan)
		}
	}()

	go func() {
		for update := range b.holidayCommands {
			b.service.UpdateHolidayCommand(update.Message, update.Message.Chat.ID, b.userState, msgChan)
		}
	}()

	go func() {
		for update := range b.weatherCommands {
			b.service.UpdateWeatherCommand(update.Message, update.Message.Chat.ID, b.userState, msgChan)
		}
	}()

	go func() {
		for update := range b.subscribeCommands {
			b.service.UpdateSubscribeCommand(&update, b.userState, msgChan)
		}
	}()

	go func() {
		for update := range b.unsubscribeCommands {
			b.service.UpdateUnsubscribeCommand(&update, b.userState, msgChan)
		}
	}()

	go func() {
		for msg := range msgChan {
			_, err := b.api.Send(msg)
			if err != nil {
				b.log.Errorf("error sending message, err: %v", err)
				return
			}
		}
	}()
}

func getChatID(update tgbotapi.Update) int64 {
	var chatID int64
	if update.Message != nil {
		chatID = update.Message.Chat.ID
	} else {
		chatID = update.CallbackQuery.From.ID
	}

	return chatID
}

func (b *Bot) createMenu() {
	cfg := tgbotapi.NewSetMyCommands(
		tgbotapi.BotCommand{
			Command:     handler.StartMenu,
			Description: "Show menu",
		},
		tgbotapi.BotCommand{
			Command:     handler.HolidayMenu,
			Description: "Show today holiday",
		},
		tgbotapi.BotCommand{
			Command:     handler.WeatherMenu,
			Description: "Show current weather",
		},
		tgbotapi.BotCommand{
			Command:     handler.SubscribeMenu,
			Description: "Subscribe to weather forecast",
		},
		tgbotapi.BotCommand{
			Command:     handler.UnsubscribeMenu,
			Description: "Unsubscribe from weather forecast",
		})

	_, _ = b.api.Request(cfg)
}
