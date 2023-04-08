// based on swag example (MIT License): https://github.com/swaggo/swag

package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
func (c *Controller) Article(ctx *gin.Context) {
	id, err := uuid.FromString(ctx.Param("id"))
	if err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, errors.New("invalid id"))
		return
	}

	article, err := c.repository.GetArticle(id)
	if err != nil {
		httputil.NewError(ctx, http.StatusNotFound, errors.New("article not found"))
		return
	}
	ctx.JSON(http.StatusOK, article)
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
func (c *Controller) Articles(ctx *gin.Context) {
	feedID, err := uuid.FromString(ctx.Param("feed_id"))
	if err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, errors.New("invalid feed_id"))
		return
	}

	limit := c.articlesPerPage
	offset, err := strconv.Atoi(ctx.Query("start"))
	if err != nil {
		offset = 0
	}

	articles, err := c.repository.GetArticles(feedID, limit, offset)
	if err != nil {
		httputil.NewError(ctx, http.StatusNotFound, errors.New("articles not found"))
		return
	}
	ctx.JSON(http.StatusOK, articles)
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
func (c *Controller) Folders(ctx *gin.Context) {
	claims := helpers.GetAuthClaims(ctx)
	folders, err := c.repository.Folders(claims)
	if err != nil {
		httputil.NewError(ctx, http.StatusNotFound, errors.New("folders not found"))
		return
	}
	ctx.JSON(http.StatusOK, folders)
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
func (c *Controller) GetFeed(ctx *gin.Context) {
	feedID, err := uuid.FromString(ctx.Param("feed_id"))
	if err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, errors.New("invalid feed_id"))
		return
	}

	feed, err := c.repository.GetFeed(feedID)
	if err != nil {
		httputil.NewError(ctx, http.StatusNotFound, errors.New("feed not found"))
		return
	}
	ctx.JSON(http.StatusOK, feed)
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
func (c *Controller) Feeds(ctx *gin.Context) {
	feeds, err := c.repository.Feeds()
	if err != nil {
		httputil.NewError(ctx, http.StatusNotFound, errors.New("feeds not found"))
		return
	}
	ctx.JSON(http.StatusOK, feeds)
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
func (c *Controller) AddFeed(ctx *gin.Context) {
	var json AddFeedRequest
	if err := ctx.ShouldBindJSON(&json); err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, errors.New("malformed JSON body"))
		return
	}
	if !helpers.IsURL(json.URL) {
		httputil.NewError(ctx, http.StatusBadRequest, errors.New("not a valid URL"))
		return
	}
	if helpers.IsHTTPURL(json.URL) {
		// it's a http URL, so first try looking up if there's a https version in the db
		httpsURL := helpers.RewriteToHTTPS(json.URL)
		if feed, err := c.repository.GetFeedByURL(httpsURL); err == nil {
			// found https version
			ctx.JSON(http.StatusOK, httputil.HTTPStatus{
				Status: "ok",
				FeedID: feed.ID,
			})
			return
		}
		// try the http version next
	}
	// look if a feed with this URL already exists
	if feed, err := c.repository.GetFeedByURL(json.URL); err == nil {
		// return the existing feed id
		ctx.JSON(http.StatusOK, httputil.HTTPStatus{
			Status: "ok",
			FeedID: feed.ID,
		})
		return
	}

	// feed doesn't already exist. add it.
	feedID, err := c.repository.AddFeed(json.URL)
	if err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusOK, httputil.HTTPStatus{
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
func (c *Controller) ChangeFolders(ctx *gin.Context) {
	claims := helpers.GetAuthClaims(ctx)

	var json ChangeFoldersRequest
	if err := ctx.ShouldBindJSON(&json); err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, errors.New("malformed JSON body"))
		return
	}

	err := c.repository.ChangeFolders(claims, json.Folders)
	if err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	// send event
	c.userEventRepository.Publish(&UserEventEnvelope{
		UserID: claims.ID,
		Payload: common.UserEventMessage{
			Type: common.MessageTypeFolderUpdate,
			Data: nil,
		},
	})

	ctx.JSON(http.StatusOK, httputil.HTTPStatus{
		Status: "ok",
	})
}

// ChangeFoldersRequest is the POST body for ChangeFolders
type ChangeFoldersRequest struct {
	Folders []Folder `json:"folders"`
}
