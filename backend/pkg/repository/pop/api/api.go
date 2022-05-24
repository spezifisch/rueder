package api

import (
	"github.com/apex/log"
	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"
	"github.com/spezifisch/rueder3/backend/pkg/api/controller"
	"github.com/spezifisch/rueder3/backend/pkg/repository/pop/models"
)

// APIPopRepository internal state
type APIPopRepository struct {
	pop *pop.Connection

	folderCountLimit     int
	folderFeedCountLimit int
}

// NewAPIPopRepository returns a FeedRepository that wraps a pop DB
func NewAPIPopRepository(db string) *APIPopRepository {
	tx, err := pop.Connect(db)
	if err != nil {
		log.WithError(err).WithField("db", db).Error("couldn't connect with pop")
		return nil
	}

	return &APIPopRepository{
		pop:                  tx,
		folderCountLimit:     100,
		folderFeedCountLimit: 1000,
	}
}

// Feeds returns all feeds
func (r *APIPopRepository) Feeds() (feeds []controller.Feed, err error) {
	allFeeds := []models.Feed{}
	if err = r.pop.All(&allFeeds); err != nil {
		log.WithError(err).Error("failed fetching feeds")
		return
	}

	feeds = make([]controller.Feed, len(allFeeds))
	for i, feed := range allFeeds {
		feeds[i].ID = feed.ID
		feeds[i].URL = feed.FeedURL
		feeds[i].Title = feed.Title.String
		feeds[i].Icon = feed.Icon.String
	}
	return
}

// GetArticle returns the article with the given id
func (r *APIPopRepository) GetArticle(id uuid.UUID) (ret controller.Article, err error) {
	article := models.Article{}
	err = r.pop.Eager().Select("*, row_number() over (partition by feed_id order by seq) as feed_seq").Find(&article, id)
	if err != nil {
		log.WithError(err).Error("failed fetching article")
		return
	}

	var enclosures []controller.ArticleEnclosure = nil
	if article.Content.Enclosures != nil && len(article.Content.Enclosures) > 0 {
		enclosures = make([]controller.ArticleEnclosure, len(article.Content.Enclosures))
		for i, e := range article.Content.Enclosures {
			enclosures[i].Length = e.Length
			enclosures[i].Type = e.Type
			enclosures[i].URL = e.URL
		}
	}

	articleContent := controller.ArticleContent{
		Authors:    article.Content.Authors,
		Tags:       article.Content.Tags,
		Enclosures: enclosures,
		Text:       article.Content.Text,
	}

	ret = controller.Article{
		ID:        article.ID,
		FeedID:    article.Feed.ID,
		FeedTitle: article.Feed.Title.String,
		FeedURL:   article.Feed.FeedURL,

		Title: article.Title.String,
		Time:  article.PostedAt,
		Link:  article.Link.String,

		Thumbnail:  article.Thumbnail.String,
		Image:      article.Image.String,
		ImageTitle: article.ImageTitle.String,

		Teaser:  article.Teaser.String,
		Content: articleContent,
	}
	return
}

// GetArticles returns articles of a feed
func (r *APIPopRepository) GetArticles(feedID uuid.UUID, limit int, offset int) (articles []controller.ArticlePreview, err error) {
	// get feed
	feed := models.Feed{}
	err = r.pop.Find(&feed, feedID)
	if err != nil {
		log.WithError(err).Error("failed fetching feed")
		return
	}

	// get articles of feed
	feedArticles := models.Articles{}
	q := r.pop.Select("*, row_number() over (partition by feed_id order by seq) as feed_seq").Where("feed_id = ?", feedID)
	if offset > 0 {
		q = q.Where("seq < ?", offset)
	}
	err = q.Limit(limit).Order("articles.seq desc").All(&feedArticles)
	if err != nil {
		// no articles
		log.WithError(err).Error("no articles in feed")
		articles = make([]controller.ArticlePreview, 0)
		return
	}

	// build return value
	feedTitle := feed.Title.String
	feedIcon := feed.Icon.String
	articles = make([]controller.ArticlePreview, len(feedArticles))
	for i := 0; i < len(feedArticles); i++ {
		article := feedArticles[i]
		articles[i].ID = article.ID
		articles[i].Seq = article.Seq
		articles[i].FeedSeq = article.FeedSeq
		articles[i].FeedTitle = feedTitle
		articles[i].FeedIcon = feedIcon
		articles[i].Time = article.PostedAt
		articles[i].Title = article.Title.String
		articles[i].Teaser = article.Teaser.String
	}
	return
}

// GetFeed returns a feed with current article count
func (r *APIPopRepository) GetFeed(id uuid.UUID) (ret controller.Feed, err error) {
	feed := models.Feed{}
	err = r.pop.Find(&feed, id)
	if err != nil {
		log.WithError(err).Error("failed fetching feed")
		return
	}

	articleCount, err := r.pop.Select("id").Where("feed_id = ?", feed.ID).Count(&models.Article{})
	if err != nil {
		log.WithError(err).Error("failed counting feed articles")
		return
	}

	ret = controller.Feed{
		ID:           feed.ID,
		Title:        feed.Title.String,
		Icon:         feed.Icon.String,
		URL:          feed.FeedURL,
		SiteURL:      feed.SiteURL.String,
		ArticleCount: articleCount,
		FetcherState: &controller.FetcherState{
			Working:     feed.FetcherState.Working,
			LastSuccess: feed.FetcherState.LastSuccess,
			LastError:   feed.FetcherState.LastError,
			Message:     feed.FetcherState.Message,
		},
	}
	return
}

// AddFeed adds a new feed and returns the new feed id
func (r *APIPopRepository) AddFeed(url string) (feedID uuid.UUID, err error) {
	feed := models.Feed{
		FeedURL:     url,
		FetchDelayS: 60 * 60,
	}
	err = r.pop.Save(&feed)
	if err == nil {
		feedID = feed.ID
	}
	return
}

// GetFeedByURL returns a feed with the given URL, an error otherwise
func (r *APIPopRepository) GetFeedByURL(url string) (ret controller.Feed, err error) {
	feed := models.Feed{}
	err = r.pop.Select("id").Where("feed_url = ?", url).First(&feed)
	if err != nil {
		return
	}

	ret = controller.Feed{
		ID: feed.ID,
	}
	return
}
