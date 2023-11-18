package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"git.foxminded.ua/foxstudent106361/holiday-bot/internal/model"
	"git.foxminded.ua/foxstudent106361/holiday-bot/pkg/rest"
)

const (
	apikeyParam  = "api_key"
	countryParam = "country"
	yearParam    = "year"
	monthParam   = "month"
	dayParam     = "day"
)

var countryCodes = map[string]string{
	"🇺🇸 USA":     "US",
	"🇬🇧 UK":      "GB",
	"🇨🇦 Canada":  "CA",
	"🇫🇷 France":  "FR",
	"🇩🇪 Germany": "DE",
	"🇯🇵 Japan":   "JP",
}

type Client struct {
	httpClient rest.BaseClient
}

func New(baseURL string) *Client {
	return &Client{httpClient: rest.BaseClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}}
}

type Fetcher interface {
	GetHolidays(date time.Time, country, apiKey string) ([]model.Holiday, error)
}

func (c *Client) GetHolidays(date time.Time, country, apiKey string) ([]model.Holiday, error) {
	filters := []rest.FilterOptions{
		{
			Field:  apikeyParam,
			Values: []string{apiKey},
		},
		{
			Field:  countryParam,
			Values: []string{countryCodes[country]},
		},
		{
			Field:  yearParam,
			Values: []string{fmt.Sprint(date.Year())},
		},
		{
			Field:  monthParam,
			Values: []string{date.Format("01")},
		},
		{
			Field:  dayParam,
			Values: []string{fmt.Sprint(date.Day())},
		},
	}

	reqURL, err := c.httpClient.BuildURL("/", filters)
	if err != nil {
		return nil, fmt.Errorf("error create request url with params, err: %w", err)
	}

	req, err := http.NewRequest(http.MethodGet, reqURL, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("error create holiday request %w", err)
	}

	response, err := c.httpClient.SendRequest(req)
	if err != nil {
		return nil, fmt.Errorf("error send request %w", err)
	}

	if !response.IsOk {
		return nil, fmt.Errorf("response error, code %d, message: %s", response.StatusCode(), response.Error.Message)
	}

	respByte, err := response.ReadBody()
	if err != nil {
		return nil, fmt.Errorf("error read response body, err: %w", err)
	}

	var holidays []model.Holiday

	err = json.Unmarshal(respByte, &holidays)
	if err != nil {
		return nil, fmt.Errorf("error decode response body, err: %w", err)
	}

	return holidays, nil
}
