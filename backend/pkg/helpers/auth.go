package helpers

import (
	"fmt"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

// GetAuthClaims parses the data from the JWT
func GetAuthClaims(c *gin.Context) *AuthClaims {
	claims := jwt.ExtractClaims(c)
	if claims == nil || claims["uid"] == nil || claims["origin"] == nil || claims["sub"] == nil {
		return nil
	}
	uid, err := uuid.FromString(claims["uid"].(string))
	if err != nil {
		return nil
	}

	origin := claims["origin"].(string)
	name := claims["sub"].(string)
	originName := fmt.Sprintf("%s:%s", origin, name)

	ret := &AuthClaims{
		ID:         uid,
		Origin:     origin,
		Name:       name,
		OriginName: originName,
	}
	if !ret.IsValid() {
		return nil
	}
	return ret
}

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
