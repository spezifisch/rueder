package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SSE godoc
// @Summary Server-Side Events endpoint
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} Article
// @Failure 401 {object} httputil.HTTPError
// @Security ApiKeyAuth
// @Router /sse [get]
func (c *Controller) SSE(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "foo")
}
