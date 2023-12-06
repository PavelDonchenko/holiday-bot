package handler

import (
	"context"
	"errors"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"git.foxminded.ua/foxstudent106361/holiday-bot/internal/client"
	mocks "git.foxminded.ua/foxstudent106361/holiday-bot/internal/client/mock"
	"git.foxminded.ua/foxstudent106361/holiday-bot/internal/model"
	"git.foxminded.ua/foxstudent106361/holiday-bot/pkg/utils"
)

func TestHandler_HandleGetWeatherByCoordinate(t *testing.T) {
	forecast := model.Forecast{
		Main: struct {
			Temp      float64 `json:"temp"`
			FeelsLike float64 `json:"feels_like"`
			TempMin   float64 `json:"temp_min"`
			TempMax   float64 `json:"temp_max"`
			Pressure  int     `json:"pressure"`
			Humidity  int     `json:"humidity"`
			SeaLevel  int     `json:"sea_level"`
			GrndLevel int     `json:"grnd_level"`
		}{
			Temp:      10,
			FeelsLike: 20,
			TempMax:   20,
			TempMin:   10,
			Pressure:  90,
		},
		Name: "LVIV",
	}

	parseForecast, err := utils.ParseForecast(forecast)
	require.NoError(t, err)

	tests := []struct {
		name    string
		fetcher func(t *testing.T) client.Fetcher
		want    string
	}{
		{
			name: "ok",
			fetcher: func(t *testing.T) client.Fetcher {
				m := mocks.NewFetcher(t)
				m.On("GetForecast", "", "35.23", "43.17").Return(&forecast, nil)
				return m
			},
			want: parseForecast,
		},
		{
			name: "fetcher error",
			fetcher: func(t *testing.T) client.Fetcher {
				m := mocks.NewFetcher(t)
				m.On("GetForecast", "", "35.23", "43.17").Return(nil, errors.New("error"))
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
				Chat:     &tgbotapi.Chat{ID: 123},
				Location: &tgbotapi.Location{Latitude: 43.17, Longitude: 35.23},
			}

			msg := h.HandleGetWeatherByCoordinate(&message)
			assert.Equal(t, tt.want, msg.Text)
		})
	}
}
