package httputil

import (
	"github.com/gofrs/uuid"
)

// HTTPStatus response
type HTTPStatus struct {
	Status  string    `json:"status" example:"success"`
	Message string    `json:"message,omitempty"`
	FeedID  uuid.UUID `json:"feed_id,omitempty"`
}
