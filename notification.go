package main

import (
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

type Notification struct {
	Topic   string `json:"topic"`
	WebHook string `json:"webhook"`
}

var notifications = make(map[string][]Notification)
var notificationsMu = new(sync.Mutex)

func SubscribeHandler(ctx *gin.Context) {
	var notification Notification
	if err := ctx.ShouldBind(&notification); err != nil || notification.Topic == "" || notification.WebHook == "" {
		ctx.JSON(http.StatusBadRequest, map[string]string{"message": "missing or invalid request notification"})
		return
	}
	notificationsMu.Lock()
	defer notificationsMu.Unlock()
	notifications[notification.Topic] = append(notifications[notification.Topic], notification)
	ctx.JSON(http.StatusOK, map[string]bool{"message": true})
}

func PushNotification(topic string) {
	notificationsMu.Lock()
	defer notificationsMu.Unlock()
	notificationList, ok := notifications[topic]
	if !ok || len(notificationList) <= 0 {
		return
	}
	for _, notification := range notificationList {
		_, err := http.Post(notification.WebHook, "application/json", strings.NewReader(topic))
		if err != nil {
			log.Fatalf("notification post [%s] error: %s\n", notification.WebHook, err.Error())
		}
	}
}
