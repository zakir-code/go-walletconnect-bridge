package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"
)

func main() {
	addr := flag.String("addr", "0.0.0.0:7000", "Server listening addr")
	level := flag.String("level", "INFO", "Server log level")
	flag.Parse()
	logLevel, err := logging.LogLevel(*level)
	if err != nil {
		panic(err.Error())
	}
	leveledBackend.SetLevel(logLevel, log.Module)
	log.SetBackend(leveledBackend)

	gin.Logger()
	gin.ForceConsoleColor()
	router := gin.Default()
	{
		router.GET("/hello", HelloHandler)
		router.GET("/info", InfoHandler)
		router.POST("/subscribe", SubscribeHandler)
		router.GET("", WebSocketHandler)
	}
	panic(router.Run(*addr))
}
