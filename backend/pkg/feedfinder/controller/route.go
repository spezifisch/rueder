// based on swag example (MIT License): https://github.com/swaggo/swag

package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/spezifisch/rueder3/backend/pkg/helpers"
)

// Feedfinder godoc
// @Summary Get list of feeds from given URL that points to a HTML site
// @Tags feed
// @Accept json
// @Produce json
// @Success 200 {object} FeedFinderResponse
// @Failure 400 {object} httputil.HTTPError
// @Failure 401 {object} httputil.HTTPError
// @Failure 403 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Security ApiKeyAuth
// @Router /feedfinder [get]
func (c *Controller) Feedfinder(ctx *fiber.Ctx) error {
	claims := helpers.GetFiberAuthClaims(ctx)
	if claims == nil {
		return fiber.ErrBadRequest
	}

	var json FeedFinderRequest
	if err := ctx.BodyParser(&json); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "malformed body")
	}
	if !helpers.IsURL(json.URL) {
		return fiber.NewError(fiber.StatusBadRequest, "not a valid URL")
	}
	siteURL := json.URL

	// TODO find feeds

	result := FeedFinderResponse{
		OK:  true,
		URL: siteURL,
	}
	return ctx.JSON(result)
}
