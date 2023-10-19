package helpers

import "github.com/gofrs/uuid"

// AuthClaims is the parsed data from the JWT
type AuthClaims struct {
	ID         uuid.UUID
	Origin     string
	Name       string
	OriginName string
}

func (a AuthClaims) IsValid() bool {
	if a.Origin == "" || a.Name == "" || a.OriginName == "" || a.OriginName == ":" {
		return false
	}
	return true
}
