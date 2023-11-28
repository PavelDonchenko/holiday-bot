package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"git.foxminded.ua/foxstudent106361/holiday-bot/internal/storage"
	"git.foxminded.ua/foxstudent106361/holiday-bot/pkg/db"
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
	ctx, cancel := context.WithCancel(context.Background())

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

	logger.Info("Creating MongoDB client...")
	mongoClient, err := db.DBconnect(ctx, cfg.Mongo.URL)
	if err != nil {
		logger.Fatal("failed create mongo client", err)
	}

	botStorage := storage.NewMongo(mongoClient, cfg)

	handlers := handler.New(ctx, logger, rClient, botStorage)

	botService := service.New(handlers)

	holidayBot := bot.New(botAPI, cfg, botService, logger)

	holidayBot.Run()
	logger.Info("Bot successfully created")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	defer cancel()

	<-ctx.Done()

	logger.Info("Shutting down")
}
