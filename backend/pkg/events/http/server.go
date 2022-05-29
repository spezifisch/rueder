package http

import (
	"github.com/apex/log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/spezifisch/rueder3/backend/pkg/events/controller"
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
	s.app = fiber.New(fiber.Config{
		// trust proxy headers in dev mode
		EnableTrustedProxyCheck: !s.isDevelopmentMode,
		TrustedProxies:          s.trustedProxies,
		ProxyHeader:             fiber.HeaderXForwardedFor,
	})

	if s.isDevelopmentMode {
		// log requests
		s.app.Use(logger.New())

		// recover from panics in dev mode
		s.app.Use(recover.New(recover.Config{
			EnableStackTrace:  true,
			StackTraceHandler: recover.ConfigDefault.StackTraceHandler,
		}))

		// add CORS support because in dev mode we usually run on a different port than the frontend
		s.app.Use(cors.New(cors.Config{
			AllowOrigins:     "*",
			AllowMethods:     "GET,POST",
			AllowHeaders:     "Origin, Content-Type, Authorization",
			AllowCredentials: false, // no cookies
			ExposeHeaders:    "Content-Length",
			MaxAge:           120 * 60, // 2h, Chrome's limit
		}))
	}

	s.initAPIv1()
}

// Run starts the server
func (s *Server) Run() {
	err := s.app.Listen(s.Bind)
	if err != nil {
		log.WithError(err).Error("http server failed")
	}
}
