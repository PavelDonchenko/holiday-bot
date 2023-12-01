package storage

import (
	"context"
	"strings"

	"git.foxminded.ua/foxstudent106361/holiday-bot/config"
	"git.foxminded.ua/foxstudent106361/holiday-bot/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Storage interface {
	Save(ctx context.Context, sub model.Subscription) (string, error)
	UpdateTime(ctx context.Context, time string) error
	GetSubscriptions(ctx context.Context, chatID int64) ([]model.Subscription, error)
	GetSubscriptionByID(ctx context.Context, long, lat float64, time string) (string, error)
	DeleteSubscription(ctx context.Context, long, lat float64, time string) error
}

type Mongo struct {
	client *mongo.Collection
}

func NewMongo(client *mongo.Client, cfg config.Config) *Mongo {
	collection := client.Database(cfg.Mongo.Database).Collection(cfg.Mongo.Collection)
	return &Mongo{client: collection}
}

func (m *Mongo) Save(ctx context.Context, sub model.Subscription) (string, error) {
	_, err := m.client.InsertOne(ctx, sub)
	if err != nil {
		return "", err
	}
	return sub.ID, nil
}

func (m *Mongo) UpdateTime(ctx context.Context, time string) error {
	id := strings.Split(time, "&")[1]
	timeToSave := strings.Split(time, "&")[0]
	_, err := m.client.UpdateOne(
		ctx,
		bson.D{{"_id", id}},
		bson.D{
			{"$set", bson.D{{"notify_time", timeToSave}}}})
	if err != nil {
		return err
	}

	return nil
}

func (m *Mongo) GetSubscriptions(ctx context.Context, chatID int64) ([]model.Subscription, error) {
	cursor, err := m.client.Find(ctx, bson.D{{"chat_id", chatID}})
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

func (m *Mongo) DeleteSubscription(ctx context.Context, long, lat float64, time string) error {
	filter := bson.D{
		{"longitude", long},
		{"latitude", lat},
		{"notify_time", time},
	}

	_, err := m.client.DeleteMany(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

func (m *Mongo) GetSubscriptionByID(ctx context.Context, long, lat float64, time string) (string, error) {
	filter := bson.D{
		{"longitude", long},
		{"latitude", lat},
		{"notify_time", time},
	}

	var sub model.Subscription

	err := m.client.FindOne(ctx, filter).Decode(&sub)
	if err != nil {
		return "", err
	}

	return sub.ID, nil
}
