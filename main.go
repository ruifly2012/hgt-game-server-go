package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"server/app"
	"server/router"
	"syscall"
	"time"
)

func main() {
	// 框架初始化
	g := app.Init()
	// 初始化路由
	router.Init(g)

	server := &http.Server{
		Addr:    ":8899",
		Handler: g,
	}
	// 监听服务
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Listen:%s\n", err)
		}
	}()
	// go g.Run(":8899")

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("shutdown server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		// 优雅关服
		fmt.Println("server shutdown: ", err)
	}

	fmt.Println("server exiting...")

}
