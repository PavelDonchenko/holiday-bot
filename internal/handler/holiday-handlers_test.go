package handler

import (
	"context"
	"errors"
	"testing"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"git.foxminded.ua/foxstudent106361/holiday-bot/internal/client"
	mocks "git.foxminded.ua/foxstudent106361/holiday-bot/internal/client/mock"
	"git.foxminded.ua/foxstudent106361/holiday-bot/internal/model"
)

func TestBuildMsg(t *testing.T) {
	tests := []struct {
		name     string
		holidays []model.Holiday
		country  string
		expected string
	}{
		{
			name:     "NoHolidays",
			holidays: []model.Holiday{},
			country:  "USA",
			expected: "Country USA, doesn't have any holiday today",
		},
		{
			name: "OneHoliday",
			holidays: []model.Holiday{
				{Name: "Thanksgiving"},
			},
			country:  "USA",
			expected: "USA today holidays: \nThanksgiving\n",
		},
		{
			name: "MultipleHolidays",
			holidays: []model.Holiday{
				{Name: "Christmas"},
				{Name: "New Year"},
				{Name: "Independence Day"},
			},
			country:  "USA",
			expected: "USA today holidays: \nChristmas\nNew Year\nIndependence Day\n",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := buildMsg(tt.holidays, tt.country)
			if result != tt.expected {
				t.Errorf("Unexpected result for %s:\nExpected: %s\nActual: %s", tt.name, tt.expected, result)
			}
		})
	}
}

func TestHandler_HandleGetHolidays(t *testing.T) {
	tests := []struct {
		name    string
		fetcher func(t *testing.T) client.Fetcher
		want    string
	}{
		{
			name: "ok",
			fetcher: func(t *testing.T) client.Fetcher {
				m := mocks.NewFetcher(t)
				m.On("GetHolidays", mock.IsType(time.Time{}), "USA").Return([]model.Holiday{{Name: "Thanksgiving"}}, nil)
				return m
			},
			want: "USA today holidays: \nThanksgiving\n",
		},
		{
			name: "fetcher error",
			fetcher: func(t *testing.T) client.Fetcher {
				m := mocks.NewFetcher(t)
				m.On("GetHolidays", time.Now(), mock.Anything).Return(nil, errors.New("error"))
				return m
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				log:     logrus.New(),
				fetcher: tt.fetcher(t),
				ctx:     context.Background(),
			}

			message := tgbotapi.Message{
				Chat: &tgbotapi.Chat{ID: 123},
				Text: "USA",
			}

			msg := h.HandleGetHolidays(&message)
			assert.Equal(t, tt.want, msg.Text)
		})
	}
}
