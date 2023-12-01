package service

import (
	"fmt"
	"regexp"
	"strings"

	"git.foxminded.ua/foxstudent106361/holiday-bot/internal/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"git.foxminded.ua/foxstudent106361/holiday-bot/internal/handler"
)

type Service interface {
	UpdateRegularCommand(message *tgbotapi.Message, chatID int64, state model.State, msgChan chan tgbotapi.MessageConfig)
	UpdateHolidayCommand(message *tgbotapi.Message, chatID int64, state model.State, msgChan chan tgbotapi.MessageConfig)
	UpdateWeatherCommand(message *tgbotapi.Message, chatID int64, state model.State, msgChan chan tgbotapi.MessageConfig)
	UpdateSubscribeCommand(update *tgbotapi.Update, state model.State, msgChan chan tgbotapi.MessageConfig)
	UpdateUnsubscribeCommand(update *tgbotapi.Update, state model.State, msgChan chan tgbotapi.MessageConfig)
	UpdateModifyTimeCommand(update *tgbotapi.Update, state model.State, msgChan chan tgbotapi.MessageConfig)
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
		fmt.Println("I here, message:", message)
		msg := b.handlers.HandleSendSubscriptions(message)

		activeFlags := state[chatID]
		activeFlags.UnsubscribeActiveFlag = true
		state[chatID] = activeFlags
		fmt.Println("I here, msg:", msg)
		msgChan <- msg
	case handler.UpdateTimeMenu:
		msg := b.handlers.HandleSendSubscriptions(message)

		activeFlags := state[chatID]
		activeFlags.UpdateTimeActiveFlag = true
		state[chatID] = activeFlags

		msgChan <- msg
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
			id, err := b.handlers.HandleCreateNotification(update.Message)
			if err != nil {
				msgChan <- errMsg(update.Message.Chat.ID)
			}
			msgChan <- b.handlers.HandleShowTime(update.Message.Chat.ID, id)
		}
	}

	if update.CallbackQuery != nil {
		err := b.handlers.HandleSaveTime(update.CallbackQuery)
		if err != nil {
			msgChan <- errMsg(update.Message.Chat.ID)
		}

		state[update.CallbackQuery.From.ID] = model.ActiveFlags{SubscribeActiveFlag: false}
		msgChan <- tgbotapi.NewMessage(update.CallbackQuery.From.ID, "Subscription successfully created")
	}
}

func (b Bot) UpdateUnsubscribeCommand(update *tgbotapi.Update, state model.State, msgChan chan tgbotapi.MessageConfig) {
	if update.CallbackQuery != nil {
		err := b.handlers.HandleDeleteSub(update.CallbackQuery)
		if err != nil {
			msgChan <- errMsg(update.Message.Chat.ID)
		}

		state[update.CallbackQuery.From.ID] = model.ActiveFlags{UnsubscribeActiveFlag: false}
		msgChan <- tgbotapi.NewMessage(update.CallbackQuery.From.ID, "Subscription successfully deleted")
	}
}

func (b Bot) UpdateModifyTimeCommand(update *tgbotapi.Update, state model.State, msgChan chan tgbotapi.MessageConfig) {
	var msg tgbotapi.MessageConfig
	if isValidTimeFormat(strings.Split(update.CallbackQuery.Data, "&")[0]) {
		err := b.handlers.HandleUpdateTime(update.CallbackQuery.Data)
		if err != nil {
			msgChan <- errMsg(update.CallbackQuery.From.ID)
		}
		state[update.CallbackQuery.From.ID] = model.ActiveFlags{UpdateTimeActiveFlag: false}
		msgChan <- tgbotapi.NewMessage(update.CallbackQuery.From.ID, "Subscription successfully updated")
	} else {
		id, err := b.handlers.HandleGetSubscriptionID(update.CallbackQuery)
		if err != nil {
			msgChan <- errMsg(update.CallbackQuery.From.ID)
		}
		msg = b.handlers.HandleShowTime(update.CallbackQuery.From.ID, id)
	}

	msgChan <- msg
}

func errMsg(chatID int64) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(chatID, "error")
}

func isValidTimeFormat(input string) bool {
	pattern := `^([01]\d|2[0-3]):([0-5]\d)$`

	re := regexp.MustCompile(pattern)

	return re.MatchString(input)
}
