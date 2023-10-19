package http

// @title rueder3 Auth Backend API
// @version 1.0
// @description Auth Backend API is called internally by loginsrv

// @contact.name spezifisch
// @contact.url https://github.com/spezifisch
// @contact.email spezifisch23@proton.me

// @license.name AGPLv3
// @license.url https://www.gnu.org/licenses/agpl-3.0.en.html

// @BasePath /
// @query.collection.format multi
func (s *Server) addRoutesAuthbackend() {
	s.app.Get("/claims", s.controller.Claims)
}
