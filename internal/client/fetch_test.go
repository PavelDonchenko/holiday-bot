package client

import (
	"errors"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"

	mocks "git.foxminded.ua/foxstudent106361/holiday-bot/pkg/client/mock"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/mock"

	"git.foxminded.ua/foxstudent106361/holiday-bot/config"
	"git.foxminded.ua/foxstudent106361/holiday-bot/internal/model"
	"git.foxminded.ua/foxstudent106361/holiday-bot/pkg/client"
)

func TestClient_GetHolidays(t *testing.T) {
	now := time.Now()
	_ = godotenv.Load("../../.env")
	cfg := config.MustLoad()

	type args struct {
		date    time.Time
		country string
	}
	tests := []struct {
		name    string
		fetcher func(t *testing.T, resp *http.Response) client.Client
		args    args
		resp    *http.Response
		want    []model.Holiday
		wantErr bool
	}{
		{
			name: "OK EXIST HOLIDAYS USA",
			resp: &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(`[{"name": "Universal Children's Day"}]`)),
			},
			fetcher: func(t *testing.T, resp *http.Response) client.Client {
				m := mocks.NewClient(t)
				m.On("BuildURL", mock.Anything, mock.Anything).Return(mock.Anything, nil)
				m.On("SendRequest", mock.AnythingOfType("*http.Request")).Return(resp, nil)

				return m
			},
			args: args{
				date:    now,
				country: "🇺🇸 USA",
			},
			want: []model.Holiday{
				{Name: "Universal Children's Day"}},
			wantErr: false,
		},
		{
			name: "ERROR BUILD URL",
			fetcher: func(t *testing.T, resp *http.Response) client.Client {
				m := mocks.NewClient(t)
				m.On("BuildURL", mock.Anything, mock.Anything).Return("", errors.New("something went wrong"))

				return m
			},
			args: args{
				date:    now,
				country: "🇺🇸 USA",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "OK ERROR SEND REQUEST",
			fetcher: func(t *testing.T, resp *http.Response) client.Client {
				m := mocks.NewClient(t)
				m.On("BuildURL", mock.Anything, mock.Anything).Return(mock.Anything, nil)
				m.On("SendRequest", mock.AnythingOfType("*http.Request")).Return(nil, errors.New("something went wrong"))

				return m
			},
			args: args{
				date:    now,
				country: "🇺🇸 USA",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "OK RESPONSE NOT 200",
			resp: &http.Response{
				StatusCode: 500,
			},
			fetcher: func(t *testing.T, resp *http.Response) client.Client {
				m := mocks.NewClient(t)
				m.On("BuildURL", mock.Anything, mock.Anything).Return(mock.Anything, nil)
				m.On("SendRequest", mock.AnythingOfType("*http.Request")).Return(resp, nil)

				return m
			},
			args: args{
				date:    now,
				country: "🇺🇸 USA",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := &Client{
				httpClient: tt.fetcher(t, tt.resp),
				cfg:        cfg,
			}
			got, err := c.GetHolidays(tt.args.date, tt.args.country)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetHolidays() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetHolidays() got = %v, want %v", got, tt.want)
			}
		})
	}
}
