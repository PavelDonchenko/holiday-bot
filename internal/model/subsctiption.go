package model

import "time"

type Subscription struct {
	ID         string    `bson:"_id"`
	ChatID     int64     `bson:"chat_id"`
	Longitude  float64   `bson:"longitude"`
	Latitude   float64   `bson:"latitude"`
	NotifyTime time.Time `bson:"notify_time"`
}
