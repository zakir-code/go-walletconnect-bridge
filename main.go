package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"
)

func main() {
	logging.SetBackend(BackendFormatter)
	gin.Logger()
	gin.ForceConsoleColor()
	router := gin.Default()
	{
		router.GET("/hello", HelloHandler)
		router.GET("/info", InfoHandler)
		router.POST("/subscribe", SubscribeHandler)
		router.GET("", WebSocketHandler)
	}
	fmt.Println(router.Run(":7000"))
}
