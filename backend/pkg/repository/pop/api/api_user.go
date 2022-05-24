package api

import (
	"errors"
	"fmt"

	"github.com/apex/log"
	mapset "github.com/deckarep/golang-set"
	"github.com/gofrs/uuid"
	"github.com/spezifisch/rueder3/backend/pkg/api/controller"
	"github.com/spezifisch/rueder3/backend/pkg/helpers"
	"github.com/spezifisch/rueder3/backend/pkg/repository/pop/models"
)

// Folders returns folders
func (r *APIPopRepository) Folders(claims *helpers.AuthClaims) (ret []controller.Folder, err error) {
	if r == nil || r.pop == nil {
		err = errors.New("invalid repository")
		return
	}
	// check login
	if claims == nil || !claims.IsValid() {
		err = errors.New("invalid claims")
		return
	}

	user := models.User{}
	err = r.pop.Find(&user, claims.ID)
	if err != nil {
		err = errors.New("user doesn't exist")
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

// ChangeFolders saves the folder structure for a user
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
	// check params
	if folders == nil {
		err = errors.New("invalid folders")
		return
	}

	// get user folder list
	user := models.User{}
	err = r.pop.Find(&user, claims.ID)
	if err != nil {
		err = errors.New("user doesn't exist")
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
