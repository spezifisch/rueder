package controller

type UserEventConsumer struct {
	Channel <-chan UserEventMessage
	Close   chan<- struct{}
}

// UserEventMessage contains a message coming from backend API
type UserEventMessage struct {
	Message []byte
}
