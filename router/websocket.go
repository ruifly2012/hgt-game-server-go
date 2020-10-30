package router

import (
	"github.com/gin-gonic/gin"
	"server/websocket"
)

func Websocket(g *gin.Engine) {
	// 加载协议
	websocket.LoadProtocol()
	g.GET("/wss", websocket.Server)
}