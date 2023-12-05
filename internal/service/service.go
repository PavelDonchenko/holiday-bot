package service

import (
	"regexp"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"git.foxminded.ua/foxstudent106361/holiday-bot/internal/model"

	"git.foxminded.ua/foxstudent106361/holiday-bot/internal/handler"
)

type Service interface {
	UpdateRegularCommand(message *tgbotapi.Message, chatID int64, state model.State, msgChan chan tgbotapi.MessageConfig)
	UpdateHolidayCommand(message *tgbotapi.Message, chatID int64, state model.State, msgChan chan tgbotapi.MessageConfig)
	UpdateWeatherCommand(message *tgbotapi.Message, chatID int64, state model.State, msgChan chan tgbotapi.MessageConfig)
	UpdateSubscribeCommand(update *tgbotapi.Update, state model.State, msgChan chan tgbotapi.MessageConfig)
	UpdateUnsubscribeCommand(update *tgbotapi.Update, state model.State, msgChan chan tgbotapi.MessageConfig)
}

type Bot struct {
	handlers handler.Handlers
}

func New(handlers handler.Handlers) *Bot {
	return &Bot{handlers: handlers}
}

func (b Bot) UpdateRegularCommand(message *tgbotapi.Message, chatID int64, state model.State, msgChan chan tgbotapi.MessageConfig) {
	switch message.Text {
	case handler.StartMenu:
		msg := b.handlers.HandleStart(message)
		msgChan <- msg
	case handler.HolidayMenu:
		msg := b.handlers.HandleFlags(message)

		activeFlags := state[chatID]
		activeFlags.HolidayActiveFlag = true
		state[chatID] = activeFlags

		msgChan <- msg
	case handler.WeatherMenu:
		msg := b.handlers.HandleSendLocation(message)

		activeFlags := state[chatID]
		activeFlags.WeatherActiveFlag = true
		state[chatID] = activeFlags

		msgChan <- msg
	case handler.SubscribeMenu:
		msg := b.handlers.HandleSendLocation(message)

		activeFlags := state[chatID]
		activeFlags.SubscribeActiveFlag = true
		state[chatID] = activeFlags

		msgChan <- msg
	case handler.UnsubscribeMenu:
		msg := b.handlers.HandleSendSubscriptions(message)

		activeFlags := state[chatID]
		activeFlags.UnsubscribeActiveFlag = true
		state[chatID] = activeFlags
		msgChan <- msg
	default:
		msgChan <- tgbotapi.NewMessage(message.Chat.ID, "unknown command")
	}
}

func (b Bot) UpdateHolidayCommand(message *tgbotapi.Message, chatID int64, state model.State, msgChan chan tgbotapi.MessageConfig) {
	msgChan <- b.handlers.HandleGetHolidays(message)
	state[chatID] = model.ActiveFlags{HolidayActiveFlag: false}
}

func (b Bot) UpdateWeatherCommand(message *tgbotapi.Message, chatID int64, state model.State, msgChan chan tgbotapi.MessageConfig) {
	msg := b.handlers.HandleGetWeatherByCoordinate(message)
	msg.ParseMode = tgbotapi.ModeHTML
	msgChan <- msg
	state[chatID] = model.ActiveFlags{WeatherActiveFlag: false}
}

func (b Bot) UpdateSubscribeCommand(update *tgbotapi.Update, state model.State, msgChan chan tgbotapi.MessageConfig) {
	if update.Message != nil {
		if update.Message.Location != nil {
			_, err := b.handlers.HandleCreateNotification(update.Message)
			if err != nil {
				msgChan <- errMsg(update.Message.Chat.ID, "failed save subscription")
			}
			msgChan <- b.handlers.HandleGetTime(update.Message.Chat.ID)
		}
	}

	if update.Message.Text != "" {
		if !isValidTimeFormat(update.Message.Text) {
			err := b.handlers.HandleDeleteLastSubscription()
			if err != nil {
				msgChan <- errMsg(update.Message.Chat.ID, "failed save subscription")
			}

			msgChan <- errMsg(update.Message.Chat.ID, "wrong time format, please use correct ('08:00')")
			state[update.Message.Chat.ID] = model.ActiveFlags{SubscribeActiveFlag: false}
			return
		}

		sub, err := b.handlers.HandleGetLastSubscription()
		if err != nil {
			msgChan <- errMsg(update.Message.Chat.ID, "failed save subscription")
		}

		err = b.handlers.HandleSaveTime(update.Message.Text, sub.ID)
		if err != nil {
			msgChan <- errMsg(update.Message.Chat.ID, "failed save time")
		}

		state[update.Message.Chat.ID] = model.ActiveFlags{SubscribeActiveFlag: false}
		msgChan <- tgbotapi.NewMessage(update.Message.Chat.ID, "Subscription successfully created")
	}
}

func (b Bot) UpdateUnsubscribeCommand(update *tgbotapi.Update, state model.State, msgChan chan tgbotapi.MessageConfig) {
	if update.CallbackQuery != nil {
		err := b.handlers.HandleDeleteSub(update.CallbackQuery)
		if err != nil {
			msgChan <- errMsg(update.Message.Chat.ID, "failed delete subscription")
		}

		state[update.CallbackQuery.From.ID] = model.ActiveFlags{UnsubscribeActiveFlag: false}
		msgChan <- tgbotapi.NewMessage(update.CallbackQuery.From.ID, "Subscription successfully deleted")
	}
}

func errMsg(chatID int64, text string) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(chatID, text)
}

func isValidTimeFormat(input string) bool {
	pattern := `^([01]\d|2[0-3]):([0-5]\d)$`

	re := regexp.MustCompile(pattern)

	return re.MatchString(input)
}
