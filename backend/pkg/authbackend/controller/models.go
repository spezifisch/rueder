package controller

import (
	"github.com/gofrs/uuid"
)

// User contains user details to include in the JWT
type User struct {
	ID uuid.UUID `json:"id"`

	AuthOrigin  string `json:"auth_origin"`
	AuthSubject string `json:"auth_subject"`
}
