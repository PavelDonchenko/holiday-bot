package storage

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"git.foxminded.ua/foxstudent106361/holiday-bot/config"
	"git.foxminded.ua/foxstudent106361/holiday-bot/internal/model"
)

//go:generate mockery --name=Storage --output=mock --case=underscore
type Storage interface {
	Save(ctx context.Context, sub model.Subscription) (string, error)
	UpdateTime(ctx context.Context, time time.Time, id string) error
	GetSubscriptions(ctx context.Context, chatID int64) ([]model.Subscription, error)
	GetLastSubscription(ctx context.Context) (model.Subscription, error)
	DeleteSubscription(ctx context.Context, long, lat float64, time string) error
	DeleteLastSubscription(ctx context.Context) error
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

func (m *Mongo) UpdateTime(ctx context.Context, time time.Time, id string) error {
	_, err := m.client.UpdateOne(
		ctx,
		bson.D{{Key: "_id", Value: id}},
		bson.D{
			{Key: "$set", Value: bson.D{{Key: "notify_time", Value: time}}}})
	if err != nil {
		return err
	}

	return nil
}

func (m *Mongo) GetSubscriptions(ctx context.Context, chatID int64) ([]model.Subscription, error) {
	cursor, err := m.client.Find(ctx, bson.D{{Key: "chat_id", Value: chatID}})
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
		{Key: "longitude", Value: long},
		{Key: "latitude", Value: lat},
		{Key: "notify_time", Value: time},
	}

	_, err := m.client.DeleteMany(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

func (m *Mongo) DeleteLastSubscription(ctx context.Context) error {
	opt := options.FindOneAndDelete()
	opt.SetSort(bson.D{{Key: "created_at", Value: -1}})

	var sub model.Subscription
	if err := m.client.FindOneAndDelete(ctx, bson.D{}, opt).Decode(&sub); err != nil {
		return err
	}

	return nil
}

func (m *Mongo) GetLastSubscription(ctx context.Context) (model.Subscription, error) {
	var sub model.Subscription

	opt := options.FindOne()
	opt.SetSort(bson.D{{Key: "created_at", Value: -1}})

	if err := m.client.FindOne(ctx, bson.D{}, opt).Decode(&sub); err != nil {
		return model.Subscription{}, err
	}

	return sub, nil
}
