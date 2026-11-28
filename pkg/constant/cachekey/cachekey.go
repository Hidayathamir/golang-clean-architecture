package cachekey

import "fmt"

func UserID(id int64) string {
	return fmt.Sprintf("user:%d", id)
}

// NotificationChannel is the Redis pub/sub channel a user's real-time
// notifications are published on. The webserver subscribes to the wildcard
// NotificationChannelPattern and pushes received payloads to local sockets.
func NotificationChannel(userID int64) string {
	return fmt.Sprintf("notif:user:%d", userID)
}

const NotificationChannelPattern = "notif:user:*"
