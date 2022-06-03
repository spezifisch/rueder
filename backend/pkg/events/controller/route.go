package controller

import (
	"bufio"
	"fmt"
	"time"

	"github.com/apex/log"
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
	claims := helpers.GetFiberAuthClaims(c)
	if claims == nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	userID := claims.ID

	// based on https://github.com/gofiber/recipes/blob/73e31998b30239a9823d6ef55c01e6eade8587cf/sse/main.go
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")

	c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		// NOTE do not access anything from the fiber/fasthttp context in here (only copies like userID)
		log.WithField("userID", userID).Info("connected")

		var i int
		for {
			i++
			msg := fmt.Sprintf("%d - the time is %v", i, time.Now())
			fmt.Fprintf(w, "event: message\ndata: Message: %s\n\n", msg)
			log.Infof("sending: %v", msg)

			err := w.Flush()
			if err != nil {
				log.WithField("userID", userID).WithError(err).Info("disconnected")
				break
			}
			time.Sleep(2 * time.Second)
		}
	}))

	return nil
}
