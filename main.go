package main

import (
	"context"
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

func main() {
	addr := flag.String("addr", "0.0.0.0:7000", "Server listening addr")
	level := flag.String("level", "INFO", "Server log level")
	isHttps := flag.Bool("https", false, "Start the HTTPS")
	certFile := flag.String("cert", "cert.pem", "TLS cert file")
	keyFile := flag.String("key", "key.pem", "TLS key file")
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
		router.GET("/", WebSocketHandler)
	}
	srv := &http.Server{Addr: *addr, Handler: router}
	go func() {
		var err error
		if *isHttps {
			log.Info("Server Listen On:", "https://"+srv.Addr)
			err = srv.ListenAndServeTLS(GetFileName(*certFile), GetFileName(*keyFile))
		} else {
			log.Info("Server Listen On:", "http://"+srv.Addr)
			err = srv.ListenAndServe()
		}
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutdown Server ...")

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Info("Server exiting.")
}

func GetFileName(name string) string {
	if _, err := os.Stat(name); err == nil {
		return name
	}
	name = filepath.Join(os.Getenv("GOPATH"),
		"/src/github.com/zhcppy/go-walletconnect-bridge", name)
	if _, err := os.Stat(name); err == nil {
		return name
	}
	panic("No found file: " + name)
}
