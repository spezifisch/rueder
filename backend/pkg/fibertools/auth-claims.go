package fibertools

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v4"

	"github.com/spezifisch/rueder3/backend/pkg/helpers"
)

// GetFiberAuthClaims parses the data from the JWT in Fiber Handlers
func GetFiberAuthClaims(c *fiber.Ctx) *helpers.AuthClaims {
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

	ret := &helpers.AuthClaims{
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
