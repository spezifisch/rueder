package http

import "github.com/gofiber/fiber/v2"

// Controller is the URL handler
type Controller interface {
	Claims(ctx *fiber.Ctx) error
}
