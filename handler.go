package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	version    = "v0.0.1-beta"
	name       = "zhcppy"
	repository = "https://github.com/zhcppy/go-walletconnect-bridge"
)

func HelloHandler(ctx *gin.Context) {
	ctx.String(http.StatusOK, fmt.Sprintf("Hello World, this is Go implement WalletConnect %s", version))
}

func InfoHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, map[string]string{
		"name":       name,
		"version":    version,
		"repository": repository,
	})
}

func SubscribeHandler(ctx *gin.Context) {
	var body struct {
		Topic   string `json:"topic"`
		Webhook string `json:"webhook"`
	}
	if err := ctx.ShouldBind(body); err != nil || body.Topic == "" || body.Webhook == "" {
		ctx.JSON(http.StatusOK, map[string]string{
			"message": "missing or invalid [topic webhook] field",
		})
	}
	ctx.JSON(http.StatusOK, map[string]bool{"message": true})
}
