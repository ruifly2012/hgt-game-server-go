package app

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"server/exception"
	"strings"
	"time"
)

// 游戏服务recover
func GameServerRecover(err interface{}) {
	fmt.Println("in recover")
	var (
		unknownErr = true // 是否未知错误
	)
	go func() {
		if unknownErr {
			stack := stack(3)
			if gin.IsDebugging() {
				// 直接输出错误信息到控制台
				fmt.Println(fmt.Sprintf("[Recovery] %s panic recovered:\n\n%s\n%s",
					time.Now(), err, stack))
				Logger().Error(fmt.Sprintf("[Recovery] %s panic recovered:\n\n%s\n%s",
					time.Now(), err, stack))
			} else {
				Logger().Error(fmt.Sprintf("[Recovery] %s panic recovered:\n%s\n%s",
					time.Now(), err, stack))
			}
		}
	}()
}

// http 异常recover
func ExceptionRecover(c *gin.Context) {
	if err := recover(); err != nil {
		var (
			unknownErr = false // 是否未知错误
			brokenPipe bool // 断开连接
			message = "" // 响应信息
			code = 100 // 错误码
			status = 200 // 状态码
		)
		if ne, ok := err.(*net.OpError); ok {
			if se, ok := ne.Err.(*os.SyscallError); ok {
				if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
					brokenPipe = true
				}
			}
		}
		switch err.(type) {
		case exception.LogicException:
			exception := err.(exception.LogicException)
			message = exception.Message
			code = exception.Code
			status = exception.Status
		case exception.AuthException:
			exception := err.(exception.AuthException)
			message = exception.Message
			code = exception.Code
			status = exception.Status
		default:
			unknownErr = true
		}
		go func() {
			if brokenPipe || unknownErr {
				stack := stack(3)
				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					Logger().Error(fmt.Sprintf("%s\n%s", err, string(httpRequest)))
				} else if gin.IsDebugging() {
					headers := strings.Split(string(httpRequest), "\r\n")
					for idx, header := range headers {
						current := strings.Split(header, ":")
						if current[0] == "Authorization" {
							headers[idx] = current[0] + ": *"
						}
					}
					headersJson, _ := json.Marshal(headers)
					// 直接输出错误信息到控制台
					fmt.Println(fmt.Sprintf("[Recovery] %s panic recovered:\n%s\n%s\n%s",
						time.Now(), headersJson, err, stack))
					Logger().Error(fmt.Sprintf("[Recovery] %s panic recovered:\n%s\n%s\n%s",
						time.Now(), headersJson, err, stack))
				} else {
					Logger().Error(fmt.Sprintf("[Recovery] %s panic recovered:\n%s\n%s",
						time.Now(), err, stack))
				}
			}
		}()
		// If the connection is dead, we can't write a status to it.
		if brokenPipe {
			c.Error(err.(error)) // nolint: errcheck
			c.Abort()
		} else {
			if unknownErr {
				c.AbortWithStatus(http.StatusInternalServerError)
			} else {
				c.JSON(status, gin.H{
					"_message": message,
					"_code": code,
					"_data": nil,
				})
			}
		}
	}
}
