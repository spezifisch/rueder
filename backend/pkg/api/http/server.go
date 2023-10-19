package http

import (
	"github.com/apex/log"
	"github.com/gofiber/fiber/v2"

	"github.com/spezifisch/rueder3/backend/pkg/api/controller"
	"github.com/spezifisch/rueder3/backend/pkg/fibertools"
)

// Server is a http server
type Server struct {
	Bind string

	app               *fiber.App
	controller        *controller.Controller
	jwtSecretKey      string
	isDevelopmentMode bool
	trustedProxies    []string
}

// NewServer creates a default http backend
func NewServer(controller *controller.Controller, jwtSecretKey string, isDevelopmentMode bool, trustedProxies []string) *Server {
	if controller == nil {
		panic("controller is nil")
	}

	s := &Server{
		Bind:              ":8080",
		controller:        controller,
		jwtSecretKey:      jwtSecretKey,
		isDevelopmentMode: isDevelopmentMode,
		trustedProxies:    trustedProxies,
	}
	s.init()
	return s
}

func (s *Server) init() {
	appName := "rueder-reader"
	if s.isDevelopmentMode {
		appName += "-dev"
	}

	// distrust proxy headers only in prod mode
	enableTrustedProxyCheck := !s.isDevelopmentMode
	s.app = fibertools.NewFiberRuederApp(appName, s.isDevelopmentMode, enableTrustedProxyCheck, nil)

	// add auth middleware, all following routes require auth
	authMiddleware, err := fibertools.NewFiberAuthMiddleware(s.jwtSecretKey)
	if err != nil {
		log.WithError(err).Error("couldn't setup jwt auth middleware")
		return
	}
	s.app.Use(authMiddleware)

	// add routes
	s.addRoutesApiReaderV1()
}

// Run starts the server
func (s *Server) Run() {
	err := s.app.Listen(s.Bind)
	if err != nil {
		log.WithError(err).Fatal("http server failed")
	}
}
