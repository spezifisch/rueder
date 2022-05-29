package http

import (
	"github.com/apex/log"

	"github.com/spezifisch/rueder3/backend/internal/auth"
)

// @title rueder3 API
// @version 1.0
// @description Feed Reader API

// @contact.name spezifisch
// @contact.url https://github.com/spezifisch
// @contact.email spezifisch-go@below.fr

// @license.name AGPL-3.0-only
// @license.url https://spdx.org/licenses/AGPL-3.0-only.html

// @BasePath /api/v1
// @query.collection.format multi

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func (s *Server) initAPIv1() {
	authMiddleware, err := auth.NewAuthMiddleware(s.jwtSecretKey)
	if err != nil {
		log.WithError(err).Error("couldn't setup jwt auth middleware")
		return
	}

	v1 := s.engine.Group("/api/v1", authMiddleware.MiddlewareFunc()) /* <- this is the important part with the auth */
	{
		// not tied so the user:
		v1.GET("/article/:id", s.controller.Article)
		v1.GET("/articles/:feed_id", s.controller.Articles)
		v1.GET("/feeds", s.controller.Feeds)
		v1.GET("/feed/:feed_id", s.controller.GetFeed)
		v1.POST("/feed", s.controller.AddFeed)
		// tied to the user:
		v1.GET("/folders", s.controller.Folders)
		v1.POST("/folders", s.controller.ChangeFolders)
	}
}
