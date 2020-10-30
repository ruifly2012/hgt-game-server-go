package mongo

import (
	"go.mongodb.org/mongo-driver/mongo"
	"server/app"
)

// 对局聊天数据
func RoundChat() *mongo.Collection {
	return app.Mongo.Database("turtle_soup").Collection("round_chat")
}