package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	version    = "v0.0.1-beta"
	name       = "zhcppy"
	repository = "https://github.com/zhcppy/go-walletconnect-bridge"
)

func HealthHandler(ctx *gin.Context)  {
	ctx.Status(http.StatusNoContent)
}

func HelloHandler(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Hello! this is Go implement WalletConnect %s.", version)
}

func InfoHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, map[string]string{
		"name":       name,
		"version":    version,
		"repository": repository,
	})
}
