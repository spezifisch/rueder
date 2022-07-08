package controller

import "github.com/gofrs/uuid"

// EventRepository is for IPC notifications from the api package
type EventRepository interface {
}

// RuederRepository is the interface to the persistent database
type RuederRepository interface {
	AddFeed(url string) (feedID uuid.UUID, err error) // HACK temporary stand-in so we can't assign RedisRepository to this
}
