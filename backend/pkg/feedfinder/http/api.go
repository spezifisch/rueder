package http

// @title feedfinder API
// @version 1.0
// @description Feed Finder API

// @contact.name spezifisch
// @contact.url https://github.com/spezifisch
// @contact.email spezifisch23@proton.me

// @license.name AGPLv3
// @license.url https://www.gnu.org/licenses/agpl-3.0.en.html

// @BasePath /
// @query.collection.format multi

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func (s *Server) initAPI() {
	s.app.Get("/feedfinder", s.controller.Feedfinder)
}
