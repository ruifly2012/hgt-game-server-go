package mongo

import (
	"go.mongodb.org/mongo-driver/mongo"
	"server/app"
)

// 对局笔记数据
func RoundNote() *mongo.Collection {
	return app.Mongo.Database("turtle_soup").Collection("round_note")
}