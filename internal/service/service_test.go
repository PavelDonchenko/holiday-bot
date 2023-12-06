package service

import (
	"errors"
	"testing"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"git.foxminded.ua/foxstudent106361/holiday-bot/internal/handler"
	mocks "git.foxminded.ua/foxstudent106361/holiday-bot/internal/handler/mock"
	"git.foxminded.ua/foxstudent106361/holiday-bot/internal/model"
)

var errGeneral = errors.New("something went wrong")

func TestBot_UpdateRegularCommand(t *testing.T) {
	tests := []struct {
		name     string
		handlers func(t *testing.T) handler.Handlers
		chatID   int64
		message  *tgbotapi.Message
		wantMsg  string
	}{
		{
			name: "command start",
			handlers: func(t *testing.T) handler.Handlers {
				m := mocks.NewHandlers(t)
				m.On("HandleStart", mock.Anything).Return(tgbotapi.MessageConfig{Text: "handle start"})
				return m
			},
			message: &tgbotapi.Message{Text: handler.StartMenu},
			wantMsg: "handle start",
		},
		{
			name: "command flag",
			handlers: func(t *testing.T) handler.Handlers {
				m := mocks.NewHandlers(t)
				m.On("HandleFlags", mock.Anything).Return(tgbotapi.MessageConfig{Text: "handle flag"})
				return m
			},
			message: &tgbotapi.Message{Text: handler.HolidayMenu},
			wantMsg: "handle flag",
		},
		{
			name: "command weather",
			handlers: func(t *testing.T) handler.Handlers {
				m := mocks.NewHandlers(t)
				m.On("HandleSendLocation", mock.Anything).Return(tgbotapi.MessageConfig{Text: "handle location"})
				return m
			},
			message: &tgbotapi.Message{Text: handler.WeatherMenu},
			wantMsg: "handle location",
		},
		{
			name: "command subscribe",
			handlers: func(t *testing.T) handler.Handlers {
				m := mocks.NewHandlers(t)
				m.On("HandleSendLocation", mock.Anything).Return(tgbotapi.MessageConfig{Text: "handle location"})
				return m
			},
			message: &tgbotapi.Message{Text: handler.SubscribeMenu},
			wantMsg: "handle location",
		},
		{
			name: "command unsubscribe",
			handlers: func(t *testing.T) handler.Handlers {
				m := mocks.NewHandlers(t)
				m.On("HandleSendSubscriptions", mock.Anything).Return(tgbotapi.MessageConfig{Text: "handle unsubscribe"})
				return m
			},
			message: &tgbotapi.Message{Text: handler.UnsubscribeMenu},
			wantMsg: "handle unsubscribe",
		},
		{
			name: "unknown command",
			handlers: func(t *testing.T) handler.Handlers {
				m := mocks.NewHandlers(t)
				return m
			},
			message: &tgbotapi.Message{Chat: &tgbotapi.Chat{ID: 123}, Text: "unknown"},
			wantMsg: "unknown command",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			h := &Bot{handlers: tt.handlers(t)}

			msg := h.UpdateRegularCommand(tt.message, tt.chatID, model.State{})

			assert.Equal(t, tt.wantMsg, msg.Text)
		})
	}
}

func TestUpdateSubscribeCommand(t *testing.T) {
	tests := []struct {
		name     string
		handlers func(t *testing.T) handler.Handlers
		chatID   int64
		update   *tgbotapi.Update
		wantMsg  string
	}{
		{
			name: "create notification ok",
			handlers: func(t *testing.T) handler.Handlers {
				m := mocks.NewHandlers(t)
				m.On("HandleCreateNotification", mock.Anything).Return("notification crated", nil)
				m.On("HandleGetTime", int64(123)).Return(tgbotapi.MessageConfig{Text: "Please type the time you want to receive notification (IN '13:00' FORMAT):"}, nil)
				return m
			},
			update: &tgbotapi.Update{
				Message: &tgbotapi.Message{
					Chat:     &tgbotapi.Chat{ID: 123},
					Location: &tgbotapi.Location{Latitude: 43.17, Longitude: 35.23},
				},
			},
			wantMsg: "Please type the time you want to receive notification (IN '13:00' FORMAT):",
		},
		{
			name: "error create notification",
			handlers: func(t *testing.T) handler.Handlers {
				m := mocks.NewHandlers(t)
				m.On("HandleCreateNotification", mock.Anything).Return("", errGeneral)
				return m
			},
			update: &tgbotapi.Update{
				Message: &tgbotapi.Message{
					Chat:     &tgbotapi.Chat{ID: 123},
					Location: &tgbotapi.Location{Latitude: 43.17, Longitude: 35.23},
				},
			},
			wantMsg: "failed save subscription",
		},
		{
			name: "ok save time",
			handlers: func(t *testing.T) handler.Handlers {
				m := mocks.NewHandlers(t)
				m.On("HandleGetLastSubscription", mock.Anything).Return(model.Subscription{
					ID:         "111",
					ChatID:     123,
					Longitude:  23.23,
					Latitude:   43.43,
					NotifyTime: time.Now(),
				}, nil)
				m.On("HandleSaveTime", "12:00", "111").Return(nil)
				return m
			},
			update: &tgbotapi.Update{
				Message: &tgbotapi.Message{
					Text: "12:00",
					Chat: &tgbotapi.Chat{ID: 123},
				},
			},
			wantMsg: "Subscription successfully created",
		},
		{
			name: "validation error",
			handlers: func(t *testing.T) handler.Handlers {
				m := mocks.NewHandlers(t)
				return m
			},
			update: &tgbotapi.Update{
				Message: &tgbotapi.Message{
					Text: "wong format",
					Chat: &tgbotapi.Chat{ID: 123},
				},
			},
			wantMsg: "wrong time format, please use correct ('08:00')",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			h := &Bot{handlers: tt.handlers(t)}

			msg := h.UpdateSubscribeCommand(tt.update, model.State{})

			assert.Equal(t, tt.wantMsg, msg.Text)
		})
	}
}
