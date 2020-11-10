package mongo

import (
	"server/app"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

// 对局笔记数据
func RecordLog() *mongo.Collection {
	return app.Mongo.Database("turtle_soup_log").Collection("log_" + time.Now().Format("20060102"))
}
