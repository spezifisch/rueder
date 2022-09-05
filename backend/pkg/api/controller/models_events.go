package controller

import (
	"github.com/gofrs/uuid"

	"github.com/spezifisch/rueder3/backend/internal/common"
)

// EventEnvelope contains a message to be sent to a specific user at event API
type UserEventEnvelope struct {
	// event recipient
	UserID uuid.UUID
	// event data
	Payload common.UserEventMessage
}
