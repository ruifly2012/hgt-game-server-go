package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"server/app"
	"server/exception"
	"server/model/mongo"
	model "server/model/mysql"
)

type ExampleApi struct {}

func (u ExampleApi) Test(c *gin.Context)  {
	userId := c.Param("userId")
	var user = &model.User{UserId: userId}
	has, _ := app.DB.Get(user)
	if !has {
		exception.Logic("用户不存在")
	}

	var data = map[string]string{
		"username": user.Username,
	}

	c.Set("data", data)
}

func (u ExampleApi) Mongo(c *gin.Context) {
	value := map[string]interface{}{
			"key1": 1,
			"key2": "value2",
		}
	insertResult, err := mongo.RoundRecord().InsertOne(context.TODO(), value)
	if err != nil {
		fmt.Println(err)
	}

	c.Set("data", insertResult)
}
