package controller

type EventUserState struct {
	Channel <-chan EventMessage
	Close   chan<- struct{}
}

// EventMessage contains a message between backend API and event API
type EventMessage struct {
	Message []byte
}
