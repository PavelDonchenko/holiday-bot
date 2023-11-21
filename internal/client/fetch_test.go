package client

import (
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/joho/godotenv"

	"git.foxminded.ua/foxstudent106361/holiday-bot/config"
	"git.foxminded.ua/foxstudent106361/holiday-bot/internal/model"
	"git.foxminded.ua/foxstudent106361/holiday-bot/pkg/client"
)

func TestClient_GetHolidays(t *testing.T) {
	twoUSAHolidayDay := time.Date(2023, 11, 20, 0, 0, 0, 0, time.UTC)
	_ = godotenv.Load("../../.env")
	cfg := config.MustLoad()
	type fields struct {
		httpClient client.BaseClient
	}

	type args struct {
		date    time.Time
		country string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []model.Holiday
		wantErr bool
	}{
		{
			name: "OK EXIST HOLIDAYS USA",
			fields: fields{httpClient: client.BaseClient{
				BaseURL: "https://holidays.abstractapi.com/v1/",
				HTTPClient: &http.Client{
					Timeout: 10 * time.Second,
				}}},
			args: args{
				date:    twoUSAHolidayDay,
				country: "🇺🇸 USA",
			},
			want: []model.Holiday{
				{Name: "Universal Children's Day"},
				{Name: "Africa Industrialization Day"}},
			wantErr: false,
		},
		{
			name: "OK NO HOLIDAY",
			fields: fields{httpClient: client.BaseClient{
				BaseURL: "https://holidays.abstractapi.com/v1/",
				HTTPClient: &http.Client{
					Timeout: 10 * time.Second,
				}}},
			args: args{
				date:    twoUSAHolidayDay,
				country: "🇬🇧 UK",
			},
			want:    []model.Holiday{},
			wantErr: false,
		},
		{
			name: "ERROR WRONG URL",
			fields: fields{httpClient: client.BaseClient{
				BaseURL: "https://WRONGURL",
				HTTPClient: &http.Client{
					Timeout: 10 * time.Second,
				}}},
			args: args{
				date:    twoUSAHolidayDay,
				country: "🇺🇸 USA",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			time.Sleep(1 * time.Second)
			//t.Parallel()
			c := &Client{
				httpClient: tt.fields.httpClient,
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
