package http

import "github.com/gin-gonic/gin"

// Controller is the URL handler
type Controller interface {
	Claims(ctx *gin.Context)
}
