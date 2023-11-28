package storage

import (
	"context"

	"git.foxminded.ua/foxstudent106361/holiday-bot/config"
	"git.foxminded.ua/foxstudent106361/holiday-bot/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Storage interface {
	Save(ctx context.Context, sub model.Subscription) error
	UpdateTime(ctx context.Context, chatID int64, time string) error
}

type Mongo struct {
	client *mongo.Collection
}

func NewMongo(client *mongo.Client, cfg config.Config) *Mongo {
	collection := client.Database(cfg.Mongo.Database).Collection(cfg.Mongo.Collection)
	return &Mongo{client: collection}
}

func (m *Mongo) Save(ctx context.Context, sub model.Subscription) error {
	_, err := m.client.InsertOne(ctx, sub)
	if err != nil {
		return err
	}

	return nil
}

func (m *Mongo) UpdateTime(ctx context.Context, chatID int64, time string) error {
	_, err := m.client.UpdateOne(
		ctx,
		bson.D{{"chat_id", chatID}},
		bson.D{
			{"$set", bson.D{{"notify_time", time}}}})
	if err != nil {
		return err
	}

	return nil
}
