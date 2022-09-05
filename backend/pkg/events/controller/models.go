package controller

import "github.com/spezifisch/rueder3/backend/internal/common"

type UserEventConsumer struct {
	Channel <-chan UserEventMessage
	Close   chan<- struct{}
}

// UserEventMessage contains a message coming from backend API
type UserEventMessage struct {
	Payload common.UserEventMessage
}
