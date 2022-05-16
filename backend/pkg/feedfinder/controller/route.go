// based on swag example (MIT License): https://github.com/swaggo/swag

package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spezifisch/rueder3/backend/pkg/helpers"
	"github.com/spezifisch/rueder3/backend/pkg/httputil"
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
func (c *Controller) Feedfinder(ctx *gin.Context) {
	var json FeedFinderRequest
	if err := ctx.ShouldBindJSON(&json); err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, errors.New("malformed JSON body"))
		return
	}
	if !helpers.IsURL(json.URL) {
		httputil.NewError(ctx, http.StatusBadRequest, errors.New("not a valid URL"))
		return
	}
	siteURL := json.URL

	// TODO find feeds

	result := FeedFinderResponse{
		OK:  true,
		URL: siteURL,
	}
	ctx.JSON(http.StatusOK, result)
}
