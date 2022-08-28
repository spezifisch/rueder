package http

import (
	"github.com/apex/log"
	"github.com/gin-gonic/gin"
)

// Server is a http server
type Server struct {
	Bind string

	engine     *gin.Engine
	controller Controller
}

// NewServer creates a default http backend
func NewServer(controller Controller, isDevelopmentMode bool) *Server {
	if !isDevelopmentMode {
		gin.SetMode("release")
	}

	s := &Server{
		Bind:       ":8080",
		engine:     gin.Default(),
		controller: controller,
	}
	s.init()
	return s
}

// @title rueder3 Auth Backend API
// @version 1.0
// @description Auth Backend API

// @contact.name spezifisch
// @contact.url https://github.com/spezifisch
// @contact.email spezifisch23@proton.me

// @license.name GPLv3
// @license.url https://www.gnu.org/licenses/gpl-3.0.en.html

// @BasePath /
// @query.collection.format multi
func (s *Server) init() {
	// never trust any proxy because this service should only be used internally by loginsrv
	if err := s.engine.SetTrustedProxies(nil); err != nil {
		panic("failed setting trusted proxies to nil")
	}

	v1 := s.engine.Group("/")
	{
		v1.GET("/claims", s.controller.Claims)
	}
}

// Run starts the server
func (s *Server) Run() {
	err := s.engine.Run(s.Bind)
	if err != nil {
		log.WithError(err).Error("gin failed")
	}
}
