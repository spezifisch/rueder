package http

import (
	"fmt"

	"github.com/apex/log"

	"github.com/gofiber/fiber/v2"

	"github.com/spezifisch/rueder3/backend/internal/auth"
)

func (s *Server) initAPIv1() {
	// add auth middleware first so all routes require auth
	authMiddleware, err := auth.NewFiberAuthMiddleware(s.jwtSecretKey)
	if err != nil {
		log.WithError(err).Error("couldn't setup jwt auth middleware")
		return
	}
	s.app.Use(authMiddleware)

	s.app.Get("/*", func(c *fiber.Ctx) error {
		msg := fmt.Sprintf("✋ %s", c.Params("*"))
		return c.SendString(msg) // => ✋ register
	})
}
