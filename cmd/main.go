package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/joho/godotenv/autoload"
	"github.com/sirupsen/logrus"

	"git.foxminded.ua/foxstudent106361/holiday-bot/config"
	"git.foxminded.ua/foxstudent106361/holiday-bot/internal/bot"
	"git.foxminded.ua/foxstudent106361/holiday-bot/internal/client"
	"git.foxminded.ua/foxstudent106361/holiday-bot/internal/handler"
	"git.foxminded.ua/foxstudent106361/holiday-bot/internal/service"
)

func main() {
	logger := logrus.New()

	cfg := config.MustLoad()
	logger.Debug("config: ", cfg)

	botAPI, err := tgbotapi.NewBotAPI(cfg.Telegram.TelegramBotToken)
	if err != nil {
		logger.Fatal(err)
	}

	botAPI.Debug = cfg.Telegram.BotDebug

	logger.Infof("Authorized on account %s", botAPI.Self.UserName)

	rClient := client.New(cfg)

	handlers := handler.New(logger, rClient)

	botService := service.New(handlers)

	holidayBot := bot.New(botAPI, cfg, botService, logger)

	holidayBot.Run()
}
