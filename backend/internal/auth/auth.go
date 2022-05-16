package auth

import (
	"errors"
	"time"

	"github.com/apex/log"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

const identityKey = "uid"

// NewAuthMiddleware configures a middleware that ensures the user has a valid JWT
func NewAuthMiddleware(jwtSecretKey string) (authMiddleware *jwt.GinJWTMiddleware, err error) {
	if jwtSecretKey == "secret" {
		log.Warn("USING INSECURE JWT SECRET KEY. Do this for local development only!")
	} else if len(jwtSecretKey) < 16 {
		err = errors.New("your JWT secret key is less than 16 characters long, we won't allow that")
		return
	}

	authMiddleware, err = jwt.New(&jwt.GinJWTMiddleware{
		SigningAlgorithm: "HS512",
		Realm:            "rueder3",
		Key:              []byte(jwtSecretKey),
		IdentityKey:      identityKey,
		Authenticator:    nil, /* not using login handler */
		PayloadFunc:      nil, /* not using login handler */
		IdentityHandler:  nil, /* use default identity handler */
		Authorizator:     nil, /* no special checks needed */
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})

	if err != nil {
		return
	}

	// When you use jwt.New(), the function is already automatically called for checking,
	// which means you don't need to call it again.
	err = authMiddleware.MiddlewareInit()
	if err != nil {
		return
	}

	return
}
