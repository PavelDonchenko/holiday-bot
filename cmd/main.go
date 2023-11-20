package main

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"

	"git.foxminded.ua/foxstudent106361/holiday-bot/internal/handler"
	"git.foxminded.ua/foxstudent106361/holiday-bot/internal/service"
	"git.foxminded.ua/foxstudent106361/holiday-bot/pkg/logging"

	"git.foxminded.ua/foxstudent106361/holiday-bot/config"
	"git.foxminded.ua/foxstudent106361/holiday-bot/internal/bot"
	"git.foxminded.ua/foxstudent106361/holiday-bot/internal/client"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("unable to load .env file: %e", err)
	}

	cfg := config.MustLoad()

	logger := logging.GetLogger()

	botAPI, err := tgbotapi.NewBotAPI(cfg.Telegram.TelegramBotToken)
	if err != nil {
		log.Fatal(err)
	}

	botAPI.Debug = cfg.Telegram.BotDebug

	logger.Infof("Authorized on account %s", botAPI.Self.UserName)

	rClient := client.New(cfg)

	handlers := handler.New(logger, rClient)

	botService := service.New(handlers)

	holidayBot := bot.New(botAPI, cfg, botService, logger)

	holidayBot.Run()
}
