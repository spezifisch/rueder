package controller

import (
	"bufio"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"

	"github.com/spezifisch/rueder3/backend/pkg/helpers"
)

// DefaultRoute godoc
// @Summary A route that echoes the JWT claims
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} dict
// @Failure 400 {object} httputil.HTTPError
// @Failure 401 {object} httputil.HTTPError
// @Security ApiKeyAuth
// @Router / [get]
func (con *Controller) DefaultRoute(c *fiber.Ctx) error {
	claims := helpers.GetFiberAuthClaims(c)
	if claims == nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	return c.JSON(fiber.Map{
		"msg":    "default route of " + c.App().Config().AppName,
		"claims": claims,
	})
}

// SSE godoc
// @Summary Server-Side Events endpoint
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} dict
// @Failure 400 {object} httputil.HTTPError
// @Failure 401 {object} httputil.HTTPError
// @Security ApiKeyAuth
// @Router /sse [get]
func (con *Controller) SSE(c *fiber.Ctx) error {
	// based on https://github.com/gofiber/recipes/blob/73e31998b30239a9823d6ef55c01e6eade8587cf/sse/main.go
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")

	c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		fmt.Println("WRITER")
		var i int
		for {
			i++
			msg := fmt.Sprintf("%d - the time is %v", i, time.Now())
			fmt.Fprintf(w, "data: Message: %s\n\n", msg)
			fmt.Println(msg)

			err := w.Flush()
			if err != nil {
				// Refreshing page in web browser will establish a new
				// SSE connection, but only (the last) one is alive, so
				// dead connections must be closed here.
				fmt.Printf("Error while flushing: %v. Closing http connection.\n", err)

				break
			}
			time.Sleep(2 * time.Second)
		}
	}))

	return nil
}
