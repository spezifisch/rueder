package http

import (
	"github.com/apex/log"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
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

// @title rueder3 Auth Backend API
// @version 1.0
// @description Auth Backend API

// @contact.name spezifisch
// @contact.url https://github.com/spezifisch
// @contact.email spezifisch23@proton.me

// @license.name AGPLv3
// @license.url https://www.gnu.org/licenses/agpl-3.0.en.html

// @BasePath /
// @query.collection.format multi
func (s *Server) init() {
	appName := "rueder-authbackend"
	if s.isDevelopmentMode {
		appName += "-dev"
	}

	s.app = fiber.New(fiber.Config{
		AppName:           appName,
		EnablePrintRoutes: false,
		// never trust any proxy because this service should only be used internally by loginsrv
		EnableTrustedProxyCheck: true,
		// enforce good behaviour by frontend
		StrictRouting: true,
		CaseSensitive: true,
	})

	// add some additional middlewares in dev mode
	if s.isDevelopmentMode {
		// log requests
		s.app.Use(logger.New())

		// recover from panics in dev mode
		s.app.Use(recover.New(recover.Config{
			EnableStackTrace:  true,
			StackTraceHandler: recover.ConfigDefault.StackTraceHandler,
		}))
	}

	s.app.Get("/claims", s.controller.Claims)
}

// Run starts the server
func (s *Server) Run() {
	err := s.app.Listen(s.Bind)
	if err != nil {
		log.WithError(err).Fatal("http server failed")
	}
}
