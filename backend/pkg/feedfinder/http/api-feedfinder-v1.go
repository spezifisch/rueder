package http

// @title feedfinder API
// @version 1.0
// @description Feed Finder API is a service that searches websites for feed URLs

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
func (s *Server) addRoutesApiFeedfinderV1() {
	s.app.Get("/feedfinder", s.controller.Feedfinder)
}
