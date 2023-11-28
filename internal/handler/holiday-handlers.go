package handler

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"
	"time"

	"git.foxminded.ua/foxstudent106361/holiday-bot/internal/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) HandleFlags(message *tgbotapi.Message) tgbotapi.MessageConfig {
	countriesKeyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🇺🇸 USA"),
			tgbotapi.NewKeyboardButton("🇬🇧 UK"),
			tgbotapi.NewKeyboardButton("🇨🇦 Canada"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🇫🇷 France"),
			tgbotapi.NewKeyboardButton("🇩🇪 Germany"),
			tgbotapi.NewKeyboardButton("🇯🇵 Japan"),
		),
	)

	msg := tgbotapi.NewMessage(message.Chat.ID, "Choose a country:")
	msg.ReplyMarkup = countriesKeyboard

	return msg
}

func (h *Handler) HandleGetHolidays(message *tgbotapi.Message) tgbotapi.MessageConfig {
	now := time.Now()

	holidays, err := h.fetcher.GetHolidays(now, message.Text)
	if err != nil {
		h.log.Error(err)
		return tgbotapi.MessageConfig{}
	}

	return tgbotapi.NewMessage(message.Chat.ID, buildMsg(holidays, message.Text))
}

func buildMsg(holidays []model.Holiday, country string) string {
	if len(holidays) < 1 {
		return fmt.Sprintf("Country %s, doesn't have any holiday today", country)
	}
	var sb strings.Builder

	msg := fmt.Sprintf("%s today holidays: \n", country)
	sb.WriteString(msg)

	for i := range holidays {
		sb.WriteString(holidays[i].Name + "\n")
	}

	return sb.String()
}

func parseForecast(forecast model.Forecast) (string, error) {
	htmlTemplate := `
		<b>Weather Forecast for {{.Name}}</b>
		Temperature: <b>{{.Main.Temp}}</b>
		Feels like: <b>{{.Main.FeelsLike}}</b>
		Min temp: <b>{{.Main.TempMin}}</b>
		Max temp: <b>{{.Main.TempMax}}</b>
		Pressure: <b>{{.Main.Pressure}}</b>
	`
	tmpl, err := template.New("weatherTemplate").Parse(htmlTemplate)
	if err != nil {
		return "", err
	}

	var tplBuffer bytes.Buffer
	if err := tmpl.Execute(&tplBuffer, forecast); err != nil {
		return "", err
	}

	return tplBuffer.String(), nil
}
