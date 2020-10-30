package mongo

import (
	"go.mongodb.org/mongo-driver/mongo"
	"server/app"
)

// 对局主数据
func RoundRecord() *mongo.Collection {
	return app.Mongo.Database("turtle_soup").Collection("round_record")
}