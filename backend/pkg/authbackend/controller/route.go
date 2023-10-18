package controller

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

// Claims godoc
// @Summary An http endpoint which provides additional information on an authenticated user.
// After successful authentication against a backend system, the endpoint gets called and the provided information
// is used to enhance the user JWT claim parameters. (Reference/Source: https://github.com/tarent/loginsrv#user-endpoint)
// @Tags loginsrv
// @Accept json
// @Produce json
// @Param origin query string true "Authentication origin (eg. simple, github, google)"
// @Param subject query string true "Authentication subject (eg. the username)"
// @Param email query string false "Email address (OAuth2 only)"
// @Param domain query string false "Domain (Google only)"
// @Param groups query string false "User groups (Gitlab only)"
// @Success 200 {object} ClaimsResponse
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /claims [get]
func (c *Controller) Claims(ctx *fiber.Ctx) error {
	authOrigin := ctx.Query("origin")
	authSubject := ctx.Query("sub")
	if authOrigin == "" || authSubject == "" {
		err := errors.New("required query param missing")
		ctx.Context().SetBodyString(err.Error())
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	user, err := c.repository.GetOrCreateUser(authOrigin, authSubject)
	if err != nil {
		ctx.Context().SetBodyString(err.Error())
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	resp := ClaimsResponse{
		Subject: user.AuthSubject,
		Origin:  user.AuthOrigin,
		UserID:  user.ID.String(),
	}
	return ctx.JSON(resp)
}

// ClaimsResponse is the returned data for loginsrv's claims endpoint,
// see https://github.com/tarent/loginsrv#user-endpoint
type ClaimsResponse struct {
	// JWT standard fields
	Subject string `json:"sub"`
	Origin  string `json:"origin"`

	// extra fields
	UserID string `json:"uid"`
}
