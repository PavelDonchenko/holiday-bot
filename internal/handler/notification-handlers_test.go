package handler

import (
	"context"
	"errors"
	"testing"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"git.foxminded.ua/foxstudent106361/holiday-bot/internal/model"
	"git.foxminded.ua/foxstudent106361/holiday-bot/internal/storage"
	mocks "git.foxminded.ua/foxstudent106361/holiday-bot/internal/storage/mock"
)

var errGeneral = errors.New("something went wrong")

func TestHandler_HandleCreateNotification(t *testing.T) {
	id := uuid.New()
	ctx := context.Background()

	tests := []struct {
		name   string
		db     func(t *testing.T) storage.Storage
		wantID string
		err    error
	}{
		{
			name: "ok",
			db: func(t *testing.T) storage.Storage {
				m := mocks.NewStorage(t)
				m.On("Save", ctx, mock.Anything).Return(id.String(), nil)
				return m
			},
			wantID: id.String(),
			err:    nil,
		},
		{
			name: "db error",
			db: func(t *testing.T) storage.Storage {
				m := mocks.NewStorage(t)
				m.On("Save", ctx, mock.Anything).Return("", errGeneral)
				return m
			},
			wantID: "",
			err:    errGeneral,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				log: logrus.New(),
				db:  tt.db(t),
				ctx: ctx,
			}

			message := tgbotapi.Message{
				Chat:     &tgbotapi.Chat{ID: 123},
				Location: &tgbotapi.Location{Latitude: 43.17, Longitude: 35.23},
			}

			id, err := h.HandleCreateNotification(&message)
			if tt.err != nil {
				assert.Equal(t, tt.err, err)
				return
			}
			assert.Equal(t, tt.wantID, id)
			assert.NoError(t, err)
		})
	}
}

func TestHandleSendSubscriptions(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		db      func(t *testing.T) storage.Storage
		wantMsg string
	}{
		{
			name: "ok",
			db: func(t *testing.T) storage.Storage {
				m := mocks.NewStorage(t)
				m.On("GetSubscriptions", ctx, mock.Anything).Return([]model.Subscription{{
					Longitude:  22.2,
					Latitude:   33.3,
					NotifyTime: time.Now(),
				}}, nil)
				return m
			},
			wantMsg: "Please choose a notification:",
		},
		{
			name: "db error",
			db: func(t *testing.T) storage.Storage {
				m := mocks.NewStorage(t)
				m.On("GetSubscriptions", ctx, mock.Anything).Return(nil, errGeneral)
				return m
			},
			wantMsg: "failed get subscriptions",
		},
		{
			name: "no subscription",
			db: func(t *testing.T) storage.Storage {
				m := mocks.NewStorage(t)
				m.On("GetSubscriptions", ctx, mock.Anything).Return([]model.Subscription{}, nil)
				return m
			},
			wantMsg: "You don't have any subscription",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				log: logrus.New(),
				db:  tt.db(t),
				ctx: ctx,
			}

			message := tgbotapi.Message{
				Chat:     &tgbotapi.Chat{ID: 123},
				Location: &tgbotapi.Location{Latitude: 43.17, Longitude: 35.23},
			}

			msg := h.HandleSendSubscriptions(&message)

			assert.Equal(t, tt.wantMsg, msg.Text)
		})
	}
}

func TestParseLocationTime(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedLon  float64
		expectedLat  float64
		expectedTime string
		expectedErr  error
	}{
		{
			name:         "ok",
			input:        "Longitude: 12.345, Latitude: 67.890, time: 12:34",
			expectedLon:  12.345,
			expectedLat:  67.890,
			expectedTime: "12:34",
			expectedErr:  nil,
		},
		{
			name:         "Invalid Input",
			input:        "Invalid Input",
			expectedLon:  0,
			expectedLat:  0,
			expectedTime: "",
			expectedErr:  errors.New("invalid input format"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			lon, lat, timeValue, err := parseLocationTime(tt.input)

			assert.Equal(t, tt.expectedLon, lon)
			assert.Equal(t, tt.expectedLat, lat)
			assert.Equal(t, tt.expectedTime, timeValue)

			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
