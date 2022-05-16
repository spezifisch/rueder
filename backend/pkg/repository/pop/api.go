package pop

import (
	"errors"
	"fmt"

	"github.com/apex/log"
	mapset "github.com/deckarep/golang-set"
	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"
	"github.com/spezifisch/rueder3/backend/pkg/api/controller"
	"github.com/spezifisch/rueder3/backend/pkg/helpers"
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

// Folders returns folders
func (r *APIPopRepository) Folders(claims *helpers.AuthClaims) (ret []controller.Folder, err error) {
	if claims == nil {
		err = errors.New("invalid claims")
		return
	}
	user := models.User{}
	err = r.pop.Find(&user, claims.ID)
	if err != nil {
		return
	}

	folders := user.Folders

	// update folders struct if we change feed titles
	updateFolders := false

	ret = make([]controller.Folder, len(folders))
	for i := 0; i < len(folders); i++ {
		ret[i].ID = folders[i].ID
		ret[i].Title = folders[i].Title
		ret[i].Feeds = make([]controller.Feed, len(folders[i].Feeds))

		for feedIdx, feed := range folders[i].Feeds {
			// take feed info from folder struct
			folderFeed := &ret[i].Feeds[feedIdx]
			folderFeed.ID = feed.ID
			folderFeed.Title = feed.Title
			folderFeed.Icon = feed.Icon

			// get updated article count and other info
			var currentFeed controller.Feed
			currentFeed, err = r.GetFeed(feed.ID)
			if err != nil {
				return
			}

			if folderFeed.Title == "" && currentFeed.Title != "" {
				log.WithField("feed_id", folderFeed.ID).Infof("old feed title '%s', new title '%s'", folderFeed.Title, currentFeed.Title)

				folderFeed.Title = currentFeed.Title
				updateFolders = true
			}

			folderFeed.Icon = currentFeed.Icon
			folderFeed.URL = currentFeed.URL
			folderFeed.SiteURL = currentFeed.SiteURL
			folderFeed.ArticleCount = currentFeed.ArticleCount
		}
	}

	if updateFolders {
		log.Info("updating folders with new titles")
		user.Folders = folders
		if updErr := r.pop.UpdateColumns(&user, "folders"); updErr != nil {
			log.WithError(updErr).Error("updating folders failed")
		}
	}

	return
}

// Labels returns labels
func (r *APIPopRepository) Labels() (ret []controller.Label, err error) {
	labels := []models.Label{}
	err = r.pop.All(&labels)
	if err != nil {
		log.WithError(err).Error("failed fetching labels")
		return
	}

	ret = make([]controller.Label, len(labels))
	for i := 0; i < len(labels); i++ {
		ret[i].ID = labels[i].ID
		ret[i].Color = labels[i].Color
		ret[i].Title = labels[i].Title
	}
	return
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

// ChangeFolders does nothing
func (r *APIPopRepository) ChangeFolders(claims *helpers.AuthClaims, folders []controller.Folder) (err error) {
	if r == nil || r.pop == nil {
		err = errors.New("invalid repository")
		return
	}

	// check login
	if claims == nil || !claims.IsValid() {
		err = errors.New("invalid claims")
		return
	}
	// get user folder list
	user := models.User{}
	err = r.pop.Find(&user, claims.ID)
	if err != nil {
		return
	}
	if len(folders) > r.folderCountLimit {
		err = errors.New("too many folders")
		return
	}

	// build a hashmap so we can ensure feeds are only subscribed to once.
	subscribedFeeds := mapset.NewSet(uuid.UUID{})
	// for comparisons against empty uuids
	nullUUID := uuid.NullUUID{}
	// the sanitized folders struct for the DB
	foldersJSON := make(models.Folders, len(folders))
	// validate folder changes and copy them
	for folderIdx, folder := range folders {
		if len(folder.Feeds) > r.folderFeedCountLimit {
			err = errors.New("too many feeds in folder")
			return
		}

		folderID := folder.ID
		if helpers.IsSameUUID(nullUUID.UUID, folderID) {
			// new folder, generate UUID
			folderID, err = uuid.NewV4()
			if err != nil {
				return
			}
		}

		foldersJSON[folderIdx].ID = folderID
		foldersJSON[folderIdx].Title = folder.Title
		foldersJSON[folderIdx].Feeds = make([]models.LightFeed, len(folder.Feeds))

		for feedIdx, feed := range folder.Feeds {
			if helpers.IsSameUUID(nullUUID.UUID, feed.ID) {
				err = errors.New("null uuid in feed")
				return
			}
			if subscribedFeeds.Contains(feed.ID) {
				err = errors.New("subscribed to a feed multiple times")
				return
			}
			foldersJSON[folderIdx].Feeds[feedIdx].ID = feed.ID

			feedInfo := models.Feed{}
			err := r.pop.Find(&feedInfo, feed.ID)
			if err != nil {
				return fmt.Errorf("feed doesn't exist: %s", feed.ID)
			}
			// this feed is good

			if feed.Title != "" {
				// allow custom title
				foldersJSON[folderIdx].Feeds[feedIdx].Title = feed.Title
			} else {
				// use default title
				foldersJSON[folderIdx].Feeds[feedIdx].Title = feedInfo.Title.String
			}
			foldersJSON[folderIdx].Feeds[feedIdx].Icon = feedInfo.Icon.String

			subscribedFeeds.Add(feed.ID)
		}
	}

	// update folders in db
	user.Folders = foldersJSON
	err = r.pop.UpdateColumns(&user, "folders")
	if err != nil {
		return
	}

	// update list of subscribed feeds to detect unsubscribed feeds (in the feed worker)
	err = r.updateUserFeeds(&user, &subscribedFeeds)
	return
}

func (r APIPopRepository) updateUserFeeds(user *models.User, feedIDs *mapset.Set) (err error) {
	if user == nil || feedIDs == nil {
		return errors.New("nil user or feedIDs")
	}

	// it turns out that update of many-to-many relationships isn't yet implemented in pop. :(
	// see: https://github.com/gobuffalo/pop/issues/136
	// therefore we do it ourself...
	userFeeds := []models.UserFeed{}
	err = r.pop.Where("user_id = ?", user.ID).All(&userFeeds)
	if err != nil {
		return
	}

	// put current feed ids into a set
	existingUserFeeds := mapset.NewSet(uuid.UUID{})
	for _, uf := range userFeeds {
		existingUserFeeds.Add(uf.FeedID)
	}

	// which feed entries need to be deleted?
	deleteUserFeeds := existingUserFeeds.Difference(*feedIDs)
	if deleteUserFeeds.Cardinality() > 0 {
		log.WithField("delete", deleteUserFeeds).WithField("count", deleteUserFeeds.Cardinality()).Info("deleting user_feeds")
		err = r.pop.RawQuery("DELETE FROM user_feeds WHERE user_id = ? AND feed_id in (?)", user.ID.String(), deleteUserFeeds.ToSlice()).Exec()
		if err != nil {
			return
		}
	}

	// which feed entries need to be added?
	addUserFeeds := (*feedIDs).Difference(existingUserFeeds)
	if addUserFeeds.Cardinality() > 0 {
		log.WithField("add", addUserFeeds).WithField("count", addUserFeeds.Cardinality()).Info("adding user_feeds")
		err = r.pop.Create(toUserFeeds(user.ID, &addUserFeeds))
		if err != nil {
			return
		}
	}

	return
}

// convert a user id with a set of feed ids to a slice of UserFeeds
func toUserFeeds(userID uuid.UUID, uf *mapset.Set) []models.UserFeed {
	if uf == nil {
		return nil
	}

	ret := make([]models.UserFeed, (*uf).Cardinality())
	for i, feedID := range (*uf).ToSlice() {
		ret[i].UserID = userID
		ret[i].FeedID = feedID.(uuid.UUID)
	}
	return ret
}
