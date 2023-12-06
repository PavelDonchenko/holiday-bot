package service

import (
	"regexp"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"git.foxminded.ua/foxstudent106361/holiday-bot/internal/handler"
	"git.foxminded.ua/foxstudent106361/holiday-bot/internal/model"
)

type Service interface {
	UpdateRegularCommand(message *tgbotapi.Message, chatID int64, state model.State) tgbotapi.MessageConfig
	UpdateHolidayCommand(message *tgbotapi.Message, chatID int64, state model.State) tgbotapi.MessageConfig
	UpdateWeatherCommand(message *tgbotapi.Message, chatID int64, state model.State) tgbotapi.MessageConfig
	UpdateSubscribeCommand(update *tgbotapi.Update, state model.State) tgbotapi.MessageConfig
	UpdateUnsubscribeCommand(update *tgbotapi.Update, state model.State) tgbotapi.MessageConfig
}

type Bot struct {
	handlers handler.Handlers
}

func New(handlers handler.Handlers) *Bot {
	return &Bot{handlers: handlers}
}

func (b Bot) UpdateRegularCommand(message *tgbotapi.Message, chatID int64, state model.State) tgbotapi.MessageConfig {
	switch message.Text {
	case handler.StartMenu:
		msg := b.handlers.HandleStart(message)
		return msg
	case handler.HolidayMenu:
		msg := b.handlers.HandleFlags(message)

		activeFlags := state[chatID]
		activeFlags.HolidayActiveFlag = true
		state[chatID] = activeFlags

		return msg
	case handler.WeatherMenu:
		msg := b.handlers.HandleSendLocation(message)

		activeFlags := state[chatID]
		activeFlags.WeatherActiveFlag = true
		state[chatID] = activeFlags

		return msg
	case handler.SubscribeMenu:
		msg := b.handlers.HandleSendLocation(message)

		activeFlags := state[chatID]
		activeFlags.SubscribeActiveFlag = true
		state[chatID] = activeFlags

		return msg
	case handler.UnsubscribeMenu:
		msg := b.handlers.HandleSendSubscriptions(message)

		activeFlags := state[chatID]
		activeFlags.UnsubscribeActiveFlag = true
		state[chatID] = activeFlags
		return msg
	default:
		return tgbotapi.NewMessage(message.Chat.ID, "unknown command")
	}
}

func (b Bot) UpdateHolidayCommand(message *tgbotapi.Message, chatID int64, state model.State) tgbotapi.MessageConfig {
	state[chatID] = model.ActiveFlags{HolidayActiveFlag: false}
	return b.handlers.HandleGetHolidays(message)
}

func (b Bot) UpdateWeatherCommand(message *tgbotapi.Message, chatID int64, state model.State) tgbotapi.MessageConfig {
	msg := b.handlers.HandleGetWeatherByCoordinate(message)
	msg.ParseMode = tgbotapi.ModeHTML
	state[chatID] = model.ActiveFlags{WeatherActiveFlag: false}
	return msg
}

func (b Bot) UpdateSubscribeCommand(update *tgbotapi.Update, state model.State) tgbotapi.MessageConfig {
	if update.Message != nil {
		if update.Message.Location != nil {
			_, err := b.handlers.HandleCreateNotification(update.Message)
			if err != nil {
				return errMsg(update.Message.Chat.ID, "failed save subscription")
			}
			return b.handlers.HandleGetTime(update.Message.Chat.ID)
		}
	}

	if update.Message.Text != "" {
		if !isValidTimeFormat(update.Message.Text) {
			return errMsg(update.Message.Chat.ID, "wrong time format, please use correct ('08:00')")
		}

		sub, err := b.handlers.HandleGetLastSubscription()
		if err != nil {
			return errMsg(update.Message.Chat.ID, "failed save subscription")
		}

		if err = b.handlers.HandleSaveTime(update.Message.Text, sub); err != nil {
			return errMsg(update.Message.Chat.ID, "failed save time")
		}

		state[update.Message.Chat.ID] = model.ActiveFlags{SubscribeActiveFlag: false}
		return tgbotapi.NewMessage(update.Message.Chat.ID, "Subscription successfully created")
	}

	return tgbotapi.MessageConfig{}
}

func (b Bot) UpdateUnsubscribeCommand(update *tgbotapi.Update, state model.State) tgbotapi.MessageConfig {
	if update.CallbackQuery != nil {
		err := b.handlers.HandleDeleteSub(update.CallbackQuery)
		if err != nil {
			return errMsg(update.Message.Chat.ID, "failed delete subscription")
		}

		state[update.CallbackQuery.From.ID] = model.ActiveFlags{UnsubscribeActiveFlag: false}
		return tgbotapi.NewMessage(update.CallbackQuery.From.ID, "Subscription successfully deleted")
	}

	return tgbotapi.MessageConfig{}
}

func errMsg(chatID int64, text string) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(chatID, text)
}

func isValidTimeFormat(input string) bool {
	pattern := `^([01]\d|2[0-3]):([0-5]\d)$`

	re := regexp.MustCompile(pattern)

	return re.MatchString(input)
}
