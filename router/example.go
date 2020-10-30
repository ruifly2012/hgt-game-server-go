package router

import (
	"github.com/gin-gonic/gin"
	"server/api"
)
var exampleApi = api.ExampleApi{}

func Example(g *gin.Engine) {
	g.GET("/example/test", exampleApi.Test)
	g.GET("/example/mongo", exampleApi.Mongo)
}
