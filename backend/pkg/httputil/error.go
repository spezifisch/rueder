// based on swag example (MIT License): https://github.com/swaggo/swag

package httputil

import (
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

// NewError example
func NewError(ctx *gin.Context, status int, err error) {
	er := HTTPError{
		Code:    status,
		Message: err.Error(),
	}
	ctx.JSON(status, er)
}

// NewStatus response
func NewStatus(ctx *gin.Context, httpStatus int, statusMessage string) {
	data := HTTPStatus{
		Status: statusMessage,
	}
	ctx.JSON(httpStatus, data)
}

// HTTPError response
type HTTPError struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"status bad request"`
}

// HTTPStatus response
type HTTPStatus struct {
	Status  string    `json:"status" example:"success"`
	Message string    `json:"message,omitempty"`
	FeedID  uuid.UUID `json:"feed_id,omitempty"`
}
