package http

import (
	"github.com/apex/log"
	"github.com/gofiber/fiber/v2"

	"github.com/spezifisch/rueder3/backend/pkg/fibertools"
)

// Server is a http server
type Server struct {
	Bind string

	app               *fiber.App
	controller        Controller
	isDevelopmentMode bool
}

// NewServer creates a default http backend
func NewServer(controller Controller, bind string, isDevelopmentMode bool) *Server {
	s := &Server{
		Bind:              bind,
		controller:        controller,
		isDevelopmentMode: isDevelopmentMode,
	}
	s.init()
	return s
}

func (s *Server) init() {
	appName := "rueder-authbackend"
	if s.isDevelopmentMode {
		appName += "-dev"
	}

	// never trust any proxy because this service should only be used internally by loginsrv
	enableTrustedProxyCheck := true
	s.app = fibertools.NewFiberRuederApp(appName, s.isDevelopmentMode, enableTrustedProxyCheck, nil)

	// add routes
	s.addRoutesAuthbackend()
}

// Run starts the server
func (s *Server) Run() {
	err := s.app.Listen(s.Bind)
	if err != nil {
		log.WithError(err).Fatal("http server failed")
	}
}
