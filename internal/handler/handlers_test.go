package handler

import (
	"testing"

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
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := buildMsg(tt.holidays, tt.country)
			if result != tt.expected {
				t.Errorf("Unexpected result for %s:\nExpected: %s\nActual: %s", tt.name, tt.expected, result)
			}
		})
	}
}
