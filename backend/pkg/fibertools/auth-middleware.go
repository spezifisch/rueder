package fibertools

import (
	"errors"

	"github.com/apex/log"
	"github.com/gofiber/fiber/v2"

	jwtware "github.com/gofiber/jwt/v3"
)

// NewFiberAuthMiddleware configures a middleware that ensures the user has a valid JWT
func NewFiberAuthMiddleware(jwtSecretKey string) (authMiddleware fiber.Handler, err error) {
	if jwtSecretKey == "secret" {
		// (2x space after emoji looks better)
		log.Warn("⚠️  USING INSECURE JWT SECRET KEY. Do this for local development only!")
	} else if len(jwtSecretKey) < 16 {
		err = errors.New("your JWT secret key is less than 16 characters long, we won't allow that")
		return
	}

	config := jwtware.Config{
		SigningKey: []byte(jwtSecretKey),

		// Signing method, used to check token signing method.
		// Optional. Default: "HS256".
		// Possible values: "HS256", "HS384", "HS512", "ES256", "ES384", "ES512", "RS256", "RS384", "RS512"
		SigningMethod: "HS512",

		// ErrorHandler defines a function which is executed for an invalid token.
		// It may be used to define a custom JWT error.
		// Optional. Default: 401 Invalid or expired JWT
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			if err.Error() == "Missing or malformed JWT" {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"code":    fiber.StatusBadRequest,
					"message": "Missing or malformed JWT",
				})
			}
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"code":    fiber.StatusUnauthorized,
				"message": "Invalid or expired JWT " + err.Error(),
			})
		},

		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "param:<name>"
		// - "cookie:<name>"
		TokenLookup: "header:Authorization",

		// AuthScheme to be used in the Authorization header.
		// Optional. Default: "Bearer".
		AuthScheme: "Bearer",
	}
	authMiddleware = jwtware.New(config)
	return
}
