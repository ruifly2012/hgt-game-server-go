package app

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var Mongo *mongo.Client

func LoadMongo() {
	Mongo = SetConnect()
}

// 连接设置
func SetConnect() *mongo.Client {
	uri := "mongodb://" + MongoConfig.Username + ":" + MongoConfig.Password + "@" + MongoConfig.Ip + ":" + MongoConfig.Port
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri).SetMaxPoolSize(20)) // 连接池
	if err != nil {
		fmt.Println(err)
	}

	return client
}
