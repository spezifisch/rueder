package http

import (
	"github.com/apex/log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/spezifisch/rueder3/backend/internal/auth"
	"github.com/spezifisch/rueder3/backend/pkg/events/controller"
	"github.com/spezifisch/rueder3/backend/pkg/helpers"
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
	appName := "rueder-events"
	if s.isDevelopmentMode {
		appName += "-dev"
	}

	s.app = fiber.New(fiber.Config{
		AppName: appName,
		// print routes in dev mode
		EnablePrintRoutes: s.isDevelopmentMode,
		// distrust proxy headers only in prod mode
		EnableTrustedProxyCheck: !s.isDevelopmentMode,
		TrustedProxies:          s.trustedProxies,
		ProxyHeader:             fiber.HeaderXForwardedFor,
		// enforce good behaviour by frontend
		StrictRouting: true,
		// why is this not true by default?
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

	// add auth middleware, all following routes require auth
	authMiddleware, err := auth.NewFiberAuthMiddleware(s.jwtSecretKey)
	if err != nil {
		log.WithError(err).Error("couldn't setup jwt auth middleware")
		return
	}
	s.app.Use(authMiddleware)

	// routes
	s.app.Get("/", func(c *fiber.Ctx) error {
		claims := helpers.GetFiberAuthClaims(c)
		if claims == nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		return c.JSON(fiber.Map{
			"msg":    "default route of " + appName,
			"claims": claims,
		})
	})
}

// Run starts the server
func (s *Server) Run() {
	err := s.app.Listen(s.Bind)
	if err != nil {
		log.WithError(err).Error("http server failed")
	}
}
