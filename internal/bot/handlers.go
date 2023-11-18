package bot

import (
	"fmt"
	"log"
	"strings"
	"time"

	"git.foxminded.ua/foxstudent106361/holiday-bot/internal/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleFlags(message *tgbotapi.Message) {
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

	_, err := b.api.Send(msg)
	if err != nil {
		log.Printf("error sending mesage, err: %v", err)
		return
	}
}

func (b *Bot) handleGetHolidays(message *tgbotapi.Message) {
	now := time.Now()

	holidays, err := b.fetcher.GetHolidays(now, message.Text, b.cfg.AbstractAPIKey)
	if err != nil {
		log.Print(err)
		return
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, buildMsg(holidays, message.Text))
	_, err = b.api.Send(msg)
	if err != nil {
		log.Println("Error sending message:", err)
	}
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
