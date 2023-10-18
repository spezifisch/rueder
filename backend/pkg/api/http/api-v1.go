package http

// @title rueder3 API
// @version 1.0
// @description Feed Reader API

// @contact.name spezifisch
// @contact.url https://github.com/spezifisch
// @contact.email spezifisch23@proton.me

// @license.name AGPL-3.0-only
// @license.url https://spdx.org/licenses/AGPL-3.0-only.html

// @BasePath /api/v1
// @query.collection.format multi

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func (s *Server) initAPIv1() {
	v1 := s.app.Group("/api/v1")
	{
		// not tied so the user:
		v1.Get("/article/:id", s.controller.Article)
		v1.Get("/articles/:feed_id", s.controller.Articles)
		v1.Get("/feeds", s.controller.Feeds)
		v1.Get("/feed/:feed_id", s.controller.GetFeed)
		v1.Post("/feed", s.controller.AddFeed)
		// tied to the user:
		v1.Get("/folders", s.controller.Folders)
		v1.Post("/folders", s.controller.ChangeFolders)
	}
}
