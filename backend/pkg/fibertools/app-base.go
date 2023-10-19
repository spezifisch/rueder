package fibertools

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func NewFiberRuederApp(appName string, isDevelopmentMode, enableTrustedProxyCheck bool, trustedProxies []string) (app *fiber.App) {
	app = fiber.New(fiber.Config{
		AppName:                 appName,
		EnablePrintRoutes:       false,
		EnableTrustedProxyCheck: enableTrustedProxyCheck,
		TrustedProxies:          trustedProxies,
		ProxyHeader:             fiber.HeaderXForwardedFor,
		// enforce good behaviour by frontend
		StrictRouting: true,
		CaseSensitive: true,
	})

	// add some additional middlewares in dev mode
	if isDevelopmentMode {
		// log requests
		app.Use(logger.New())

		// recover from panics in dev mode
		app.Use(recover.New(recover.Config{
			EnableStackTrace:  true,
			StackTraceHandler: recover.ConfigDefault.StackTraceHandler,
		}))

		// add CORS support because in dev mode we usually run on a different port than the frontend
		app.Use(cors.New(cors.Config{
			AllowOrigins:     "*",
			AllowMethods:     "GET",
			AllowHeaders:     "Origin, Content-Type, Authorization",
			AllowCredentials: false, // no cookies
			ExposeHeaders:    "Content-Length",
			MaxAge:           120 * 60, // 2h, Chrome's limit
		}))
	}

	return
}
