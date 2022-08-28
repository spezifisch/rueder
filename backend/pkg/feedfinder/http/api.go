package http

import (
	"github.com/apex/log"

	"github.com/spezifisch/rueder3/backend/internal/auth"
)

// @title feedfinder API
// @version 1.0
// @description Feed Finder API

// @contact.name spezifisch
// @contact.url https://github.com/spezifisch
// @contact.email spezifisch23@proton.me

// @license.name GPLv3
// @license.url https://www.gnu.org/licenses/gpl-3.0.en.html

// @BasePath /
// @query.collection.format multi

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func (s *Server) initAPI() {
	authMiddleware, err := auth.NewAuthMiddleware(s.jwtSecretKey)
	if err != nil {
		log.WithError(err).Error("couldn't setup jwt auth middleware")
		return
	}

	v1 := s.engine.Group("/", authMiddleware.MiddlewareFunc()) /* <- this is the important part with the auth */
	{
		v1.GET("/feedfinder", s.controller.Feedfinder)
	}
}
