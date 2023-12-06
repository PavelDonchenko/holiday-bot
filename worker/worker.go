package worker

import (
	"context"
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"git.foxminded.ua/foxstudent106361/holiday-bot/config"
	"git.foxminded.ua/foxstudent106361/holiday-bot/internal/client"
	"git.foxminded.ua/foxstudent106361/holiday-bot/internal/model"
	"git.foxminded.ua/foxstudent106361/holiday-bot/pkg/utils"
)

type Worker struct {
	api     *tgbotapi.BotAPI
	db      *mongo.Client
	cfg     config.Config
	log     *logrus.Logger
	fetcher client.Fetcher
	errChan chan error
}

func New(api *tgbotapi.BotAPI, db *mongo.Client, cfg config.Config, log *logrus.Logger, fetcher client.Fetcher) *Worker {
	return &Worker{
		api:     api,
		db:      db,
		cfg:     cfg,
		log:     log,
		fetcher: fetcher,
		errChan: make(chan error),
	}
}

func (w *Worker) Run(ctx context.Context) {
	duration := time.Duration(w.cfg.Worker.WorkerDuration) * time.Second

	ticker := time.NewTicker(duration)

	for {
		select {
		case <-ticker.C:
			w.process(ctx)
		case err := <-w.errChan:
			w.log.Errorf("failed process ferecast notification, err: %v", err)
			break
		case <-ctx.Done():
			w.log.Error(ctx.Err())
			break
		}
	}
}

func (w *Worker) process(ctx context.Context) {
	w.log.Info("Starting worker processing...")
	subscriptions, err := w.getSubscriptions(ctx)
	if err != nil {
		w.errChan <- err
	}

	for i, _ := range subscriptions {
		currentTime := time.Now().UTC().Format("15:04")
		subTime := subscriptions[i].NotifyTime.Format("15:04")
		if subTime != currentTime {
			continue
		}

		forecast, err := w.fetcher.GetForecast("", fmt.Sprint(subscriptions[i].Longitude), fmt.Sprint(subscriptions[i].Latitude))
		if err != nil {
			w.errChan <- err
		}

		msg, err := utils.ParseForecast(*forecast)
		if err != nil {
			w.errChan <- err
		}

		m := tgbotapi.NewMessage(subscriptions[i].ChatID, msg)
		m.ParseMode = tgbotapi.ModeHTML

		_, err = w.api.Send(m)
		if err != nil {
			w.errChan <- err
		}
	}
}

func (w *Worker) getSubscriptions(ctx context.Context) ([]model.Subscription, error) {
	collection := w.db.Database(w.cfg.Mongo.Database).Collection(w.cfg.Mongo.Collection)
	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	var subs []model.Subscription
	for cursor.Next(ctx) {
		var sub model.Subscription
		if err = cursor.Decode(&sub); err != nil {
			return nil, err
		}
		subs = append(subs, sub)
	}

	if err = cursor.Err(); err != nil {
		return nil, err
	}

	return subs, nil
}
