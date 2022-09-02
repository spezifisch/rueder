package controller

import "github.com/gofrs/uuid"

type UserEventPublisher struct {
	Channel chan<- UserEventEnvelope
}

// EventEnvelope contains a message to be sent to a specific user at event API
type UserEventEnvelope struct {
	// event recipient
	UserID uuid.UUID
	// message
	Message []byte
}
