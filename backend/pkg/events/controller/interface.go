package controller

import "github.com/gofrs/uuid"

// UserEventRepository is for IPC notifications from the api package
type UserEventRepository interface {
	ConnectUser(uuid uuid.UUID) (state UserEventConsumer, err error)
}

// RuederRepository is the interface to the persistent database
type RuederRepository interface {
	AddFeed(url string) (feedID uuid.UUID, err error) // HACK temporary stand-in so we can't assign RedisRepository to this
}
