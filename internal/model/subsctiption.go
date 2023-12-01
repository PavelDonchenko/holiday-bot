package model

type Subscription struct {
	ID         string  `bson:"_id"`
	ChatID     int64   `bson:"chat_id"`
	Longitude  float64 `bson:"longitude"`
	Latitude   float64 `bson:"latitude"`
	NotifyTime string  `bson:"notify_time"`
}
