package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/joho/godotenv/autoload"
	"github.com/sirupsen/logrus"

	"git.foxminded.ua/foxstudent106361/holiday-bot/config"
	"git.foxminded.ua/foxstudent106361/holiday-bot/internal/bot"
	"git.foxminded.ua/foxstudent106361/holiday-bot/internal/client"
	"git.foxminded.ua/foxstudent106361/holiday-bot/internal/handler"
	"git.foxminded.ua/foxstudent106361/holiday-bot/internal/service"
	"git.foxminded.ua/foxstudent106361/holiday-bot/internal/storage"
	"git.foxminded.ua/foxstudent106361/holiday-bot/pkg/db"
	"git.foxminded.ua/foxstudent106361/holiday-bot/worker"
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

	switch cfg.Application.Run {
	case "bot":
		holidayBot := bot.New(botAPI, cfg, botService, logger)

		holidayBot.Run()
		logger.Info("Bot successfully created")
	case "worker":
		w := worker.New(botAPI, mongoClient, cfg, logger, rClient)
		logger.Info("Worker starting...")
		w.Run(ctx)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	cancel()

	<-ctx.Done()

	logger.Info("Shutting down")
}

func usage() {
	fmt.Println("Usage: [bot|worker]")
	os.Exit(1)
}
