package helpers

import (
	"fmt"

	ginJWT "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v4"
)

// GetAuthClaims parses the data from the JWT in Gin Handlers
func GetAuthClaims(c *gin.Context) *AuthClaims {
	claims := ginJWT.ExtractClaims(c)
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

// GetFiberAuthClaims parses the data from the JWT in Fiber Handlers
func GetFiberAuthClaims(c *fiber.Ctx) *AuthClaims {
	user := c.Locals("user").(*jwt.Token)
	if user == nil {
		return nil
	}
	claims := user.Claims.(jwt.MapClaims)
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
