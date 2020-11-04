package router

import (
	"github.com/gin-gonic/gin"
	"server/api"
)
var authApi = api.AuthApi{}

func Auth(g *gin.Engine) {
	// 小程序登录
	g.POST("/auth/appletLogin", authApi.AppletLogin)
	// 手机/邮箱注册
	g.POST("/auth/register", authApi.Register)
	// 账密登录
	g.POST("/auth/login", authApi.Login)
}
