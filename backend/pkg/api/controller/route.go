// SPDX-FileCopyrightText: 2022 spezifisch <spezifisch23@proton.me>
// SPDX-License-Identifier: AGPL-3.0-only
// based on swag example (MIT licensed): https://github.com/swaggo/swag

package controller

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"

	"github.com/spezifisch/rueder3/backend/internal/common"
	"github.com/spezifisch/rueder3/backend/pkg/helpers"
	"github.com/spezifisch/rueder3/backend/pkg/httputil"
)

// Article godoc
// @Summary Get article
// @Tags feed
// @Accept json
// @Produce json
// @Param id path string true "Article ID"
// @Success 200 {object} Article
// @Failure 400 {object} httputil.HTTPError
// @Failure 401 {object} httputil.HTTPError
// @Failure 403 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Security ApiKeyAuth
// @Router /article/{id} [get]
func (c *Controller) Article(ctx *fiber.Ctx) error {
	id, err := uuid.FromString(ctx.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}

	article, err := c.repository.GetArticle(id)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "article not found")
	}

	return ctx.JSON(article)
}

// Articles godoc
// @Summary Get article list
// @Tags feed
// @Accept json
// @Produce json
// @Param feed_id path  string true  "Feed ID"
// @Param start   query int    false "Start Token"
// @Success 200 {object} []ArticlePreview
// @Failure 400 {object} httputil.HTTPError
// @Failure 401 {object} httputil.HTTPError
// @Failure 403 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Security ApiKeyAuth
// @Router /articles/{feed_id} [get]
func (c *Controller) Articles(ctx *fiber.Ctx) error {
	feedID, err := uuid.FromString(ctx.Params("feed_id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid feed_id")
	}

	limit := c.articlesPerPage
	offset, err := strconv.Atoi(ctx.Query("start"))
	if err != nil {
		offset = 0
	}

	articles, err := c.repository.GetArticles(feedID, limit, offset)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "articles not found")
	}

	return ctx.JSON(articles)
}

// Folders godoc
// @Summary Get folder list
// @Tags feed
// @Accept json
// @Produce json
// @Success 200 {object} []Folder
// @Failure 400 {object} httputil.HTTPError
// @Failure 401 {object} httputil.HTTPError
// @Failure 403 {object} httputil.HTTPError
// @Security ApiKeyAuth
// @Router /folders [get]
func (c *Controller) Folders(ctx *fiber.Ctx) error {
	claims := helpers.GetFiberAuthClaims(ctx)
	folders, err := c.repository.Folders(claims)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "folders not found")
	}

	return ctx.JSON(folders)
}

// GetFeed godoc
// @Summary Get feed info for a single feed
// @Tags feed
// @Accept json
// @Produce json
// @Param feed_id path  string true  "Feed ID"
// @Success 200 {object} Feed
// @Failure 400 {object} httputil.HTTPError
// @Failure 401 {object} httputil.HTTPError
// @Failure 403 {object} httputil.HTTPError
// @Security ApiKeyAuth
// @Router /feed/{feed_id} [get]
func (c *Controller) GetFeed(ctx *fiber.Ctx) error {
	feedID, err := uuid.FromString(ctx.Params("feed_id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid feed_id")
	}

	feed, err := c.repository.GetFeed(feedID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "feed not found")
	}

	return ctx.JSON(feed)
}

// Feeds godoc
// @Summary Get all feeds
// @Tags feed
// @Accept json
// @Produce json
// @Success 200 {object} []Feed
// @Failure 400 {object} httputil.HTTPError
// @Failure 401 {object} httputil.HTTPError
// @Failure 403 {object} httputil.HTTPError
// @Security ApiKeyAuth
// @Router /feeds [get]
func (c *Controller) Feeds(ctx *fiber.Ctx) error {
	feeds, err := c.repository.Feeds()
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "feeds not found")
	}
	return ctx.JSON(feeds)
}

// AddFeed godoc
// @Summary Add feed
// @Tags feed
// @Accept json
// @Produce json
// @Param request body AddFeedRequest true "Add Feed Request"
// @Success 200 {object} httputil.HTTPStatus
// @Failure 400 {object} httputil.HTTPError
// @Failure 401 {object} httputil.HTTPError
// @Failure 403 {object} httputil.HTTPError
// @Security ApiKeyAuth
// @Router /feed [post]
func (c *Controller) AddFeed(ctx *fiber.Ctx) error {
	var json AddFeedRequest
	if err := ctx.BodyParser(&json); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "malformed JSON body")
	}
	if !helpers.IsURL(json.URL) {
		return fiber.NewError(fiber.StatusBadRequest, "not a valid URL")
	}
	if helpers.IsHTTPURL(json.URL) {
		// it's a http URL, so first try looking up if there's a https version in the db
		httpsURL := helpers.RewriteToHTTPS(json.URL)
		if feed, err := c.repository.GetFeedByURL(httpsURL); err == nil {
			// found https version
			return ctx.JSON(httputil.HTTPStatus{
				Status: "ok",
				FeedID: feed.ID,
			})
		}
		// try the http version next
	}
	// look if a feed with this URL already exists
	if feed, err := c.repository.GetFeedByURL(json.URL); err == nil {
		// return the existing feed id
		return ctx.JSON(httputil.HTTPStatus{
			Status: "ok",
			FeedID: feed.ID,
		})
	}

	// feed doesn't already exist. add it.
	feedID, err := c.repository.AddFeed(json.URL)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return ctx.JSON(httputil.HTTPStatus{
		Status: "ok",
		FeedID: feedID,
	})
}

// AddFeedRequest is the POST body for AddFeed
type AddFeedRequest struct {
	URL string `json:"url"`
}

// ChangeFolders godoc
// @Summary Change folder list, titles, feeds
// @Tags feed
// @Accept json
// @Produce json
// @Param request body ChangeFoldersRequest true "Change Folders"
// @Success 200 {object} httputil.HTTPStatus
// @Failure 400 {object} httputil.HTTPError
// @Failure 401 {object} httputil.HTTPError
// @Failure 403 {object} httputil.HTTPError
// @Security ApiKeyAuth
// @Router /folders [post]
func (c *Controller) ChangeFolders(ctx *fiber.Ctx) error {
	claims := helpers.GetFiberAuthClaims(ctx)

	var json ChangeFoldersRequest
	if err := ctx.BodyParser(&json); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "malformed JSON body")
	}

	err := c.repository.ChangeFolders(claims, json.Folders)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// send event
	c.userEventRepository.Publish(&UserEventEnvelope{
		UserID: claims.ID,
		Payload: common.UserEventMessage{
			Type: common.MessageTypeFolderUpdate,
			Data: nil,
		},
	})

	return ctx.JSON(httputil.HTTPStatus{
		Status: "ok",
	})
}

// ChangeFoldersRequest is the POST body for ChangeFolders
type ChangeFoldersRequest struct {
	Folders []Folder `json:"folders"`
}
