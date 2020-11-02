package app

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// 全局返回 json
func ReturnJson() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		data, _ :=  c.Get("data")
		if data != nil {
			c.JSON(http.StatusOK, gin.H{
				"_message": "success",
				"_code": 0,
				"_data": data,
			})
		}
	}
}

// 异常捕获恢复程序中间件
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer ExceptionRecover(c)
		c.Next()
	}
}

// 请求日志
func LoggerToFile() gin.HandlerFunc {
	logger := Logger()
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()

		// 执行时间
		latencyTime := endTime.Sub(startTime)

		// 请求方式
		reqMethod := c.Request.Method

		// 请求路由
		reqUri := c.Request.RequestURI

		// 状态码
		statusCode := c.Writer.Status()

		// 请求IP
		clientIP := c.ClientIP()

		//日志格式
		logger.Infof("| %3d | %13v | %15s | %s | %s |",
			statusCode,
			latencyTime,
			clientIP,
			reqMethod,
			reqUri,
		)
	}
}