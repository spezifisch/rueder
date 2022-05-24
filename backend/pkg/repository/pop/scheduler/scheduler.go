package scheduler

import (
	"context"
	"encoding/json"
	"time"

	"github.com/apex/log"
	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"

	"github.com/spezifisch/rueder3/backend/pkg/helpers"
	"github.com/spezifisch/rueder3/backend/pkg/repository/pop/models"
	"github.com/spezifisch/rueder3/backend/pkg/worker/scheduler"
)

// SchedulerPopRepository internal state
type SchedulerPopRepository struct {
	pop *pop.Connection
	pgx *pgx.Conn
}

// NewSchedulerPopRepository returns a SchedulerRepository that wraps a pop DB
func NewSchedulerPopRepository(db string) *SchedulerPopRepository {
	// connect using pop for ORM stuff
	popTx, err := pop.Connect(db)
	if err != nil {
		log.WithError(err).WithField("db", db).Error("couldn't connect with pop")
		return nil
	}

	return &SchedulerPopRepository{
		pop: popTx,
		pgx: nil, // will connect later
	}
}

// Feeds returns the list of feeds to fetch next for the scheduler
func (r *SchedulerPopRepository) Feeds() (ret []scheduler.Feed, err error) {
	feeds := models.Feeds{}
	err = r.pop.All(&feeds)
	if err != nil {
		return
	}

	ret = make([]scheduler.Feed, len(feeds))
	for i, feed := range feeds {
		ret[i] = r.toSchedulerFeed(&feed)
	}

	return
}

// GetFeed returns the feed with the given id
func (r *SchedulerPopRepository) GetFeed(feedID uuid.UUID) (wanted scheduler.Feed, err error) {
	feed := models.Feed{}
	err = r.pop.Find(&feed, feedID)
	if err != nil {
		return
	}

	wanted = r.toSchedulerFeed(&feed)
	return
}

func (r *SchedulerPopRepository) toSchedulerFeed(feed *models.Feed) scheduler.Feed {
	articleCount, _ := r.pop.Where("feed_id = ?", feed.ID).Select("0 as feed_seq").Count(&models.Article{})

	return scheduler.Feed{
		ID:        feed.ID,
		CreatedAt: feed.CreatedAt,
		Title:     feed.Title.String,
		FeedURL:   feed.FeedURL,
		SiteURL:   feed.SiteURL.String,
		Icon:      feed.Icon.String,
		FetcherState: scheduler.FeedFetcherState{
			FetchedAt:   feed.FetchedAt,
			FetchDelayS: feed.FetchDelayS,

			ETag:         feed.FetcherState.ETag,
			LastModified: feed.FetcherState.LastModified,

			Working:     feed.FetcherState.Working,
			LastSuccess: feed.FetcherState.LastSuccess,
			LastError:   feed.FetcherState.LastError,
			Message:     feed.FetcherState.Message,
		},
		ArticleCount: articleCount,
	}
}

// UpdateFeedInfo updates metadata and sets the fetched timestamp for the feed
func (r *SchedulerPopRepository) UpdateFeedInfo(feedID uuid.UUID, updatedFeed *scheduler.Feed) (err error) {
	feed := models.Feed{
		ID:          feedID,
		FetchedAt:   updatedFeed.FetcherState.FetchedAt,
		FetchDelayS: updatedFeed.FetcherState.FetchDelayS,
		FetcherState: models.FetcherState{
			ETag:         updatedFeed.FetcherState.ETag,
			LastModified: updatedFeed.FetcherState.LastModified,

			Working:     updatedFeed.FetcherState.Working,
			LastSuccess: updatedFeed.FetcherState.LastSuccess,
			LastError:   updatedFeed.FetcherState.LastError,
			Message:     updatedFeed.FetcherState.Message,
		},
		FeedURL: updatedFeed.FeedURL,
		SiteURL: helpers.NullStringify(updatedFeed.SiteURL),
		Title:   helpers.NullStringify(updatedFeed.Title),
		Icon:    helpers.NullStringify(updatedFeed.Icon),
	}
	err = r.pop.UpdateColumns(&feed,
		"fetched_at",
		"fetch_delay_s",
		"fetcher_state",
		"feed_url",
		"site_url",
		"title",
		"icon",
	)
	return
}

// CheckExistingArticles returns bool=true for every given SiteGUID that already exists
func (r *SchedulerPopRepository) CheckExistingArticles(feedID uuid.UUID, articleGUIDs []string) (exists []bool, err error) {
	exists = make([]bool, len(articleGUIDs))
	if len(articleGUIDs) == 0 {
		return
	}

	// build a map for better random access: guid => index in exists array
	guidIndices := make(map[string]int, len(articleGUIDs))
	for i, guid := range articleGUIDs {
		// note: the feed might be tricky here and guids might be duplicates. do we need to handle this?
		guidIndices[guid] = i

		if guid == "" {
			// this just isn't a usable guid. ignore this article
			exists[i] = true
		}
	}

	query := r.pop.Select("site_guid").Where("site_guid in (?)", articleGUIDs).Where("feed_id = ?", feedID)
	articles := []models.Article{}
	err = query.All(&articles)
	if err != nil {
		return
	}

	for _, article := range articles {
		if idx, ok := guidIndices[article.SiteGUID]; ok {
			exists[idx] = true
		}
	}

	return
}

// AddArticle stores the given article
func (r *SchedulerPopRepository) AddArticle(feedID uuid.UUID, a *scheduler.Article) (err error) {
	var enclosures []models.ArticleEnclosure = nil
	if len(a.Enclosures) > 0 {
		enclosures = make([]models.ArticleEnclosure, len(a.Enclosures))
		for i, e := range a.Enclosures {
			enclosures[i].Length = e.Length
			enclosures[i].Type = e.Type
			enclosures[i].URL = e.URL
		}
	}

	article := models.Article{
		FeedID:     feedID,
		SiteGUID:   a.SiteGUID,
		PostedAt:   a.Time.UTC(), // for some reason this doesn't get converted automatically when inserting
		Link:       helpers.NullStringify(a.Link),
		Image:      helpers.NullStringify(a.Image),
		ImageTitle: helpers.NullStringify(a.ImageTitle),
		Title:      helpers.NullStringify(a.Title),
		Teaser:     helpers.NullStringify(a.Teaser),
		Content: models.ArticleContent{
			Authors:    a.Authors,
			Tags:       a.Tags,
			Enclosures: enclosures,
			Text:       a.Text,
		},
	}

	_, err = r.pop.ValidateAndSave(&article)
	return
}

// RunFeedChangeListener adds a postgres table insert listener for the feed table
func (r *SchedulerPopRepository) RunFeedChangeListener(addedFeeds chan<- uuid.UUID, needRehash chan<- bool) (err error) {
	if r.pop.Dialect.Name() != "postgres" {
		log.WithField("dialect", r.pop.Dialect.Name()).Error("RunAddFeedListener only supports postgres. Can't use AddFeed notifications.")
		return
	}

	// connect to db
	r.connectListener()
	if r.pgx == nil {
		log.Error("pgx is nil. this shouldn't happen")
		return
	}

	// create notifier
	for _, query := range []string{
		postgresDropExistingTriggerSQL,
		postgresCreateNotifyFunctionSQL,
		postgresCreateTriggerSQL,
	} {
		_, err := r.pgx.Exec(context.Background(), query)
		if err != nil {
			return err
		}
	}

	// add listener for our channel
	_, err = r.pgx.Exec(context.Background(), "LISTEN feed_change")
	if err != nil {
		return
	}

	log.Info("added feed_change trigger and listener")

	go r.addFeedListener(addedFeeds, needRehash)
	return
}

func (r *SchedulerPopRepository) connectListener() (reconnectable bool) {
	// also connect using pgx to LISTEN for table changes
	if r.pop.Dialect.Name() == "postgres" {
		if pgxConn, err := pgx.Connect(context.Background(), r.pop.URL()); err == nil {
			r.pgx = pgxConn
		} else {
			log.WithError(err).Error("couldn't connect with pgx")
		}
	} else {
		log.Error("can't add listener for something that's not postgresql")
		return false
	}

	return true
}

func (r *SchedulerPopRepository) addFeedListener(addedFeeds chan<- uuid.UUID, needRehash chan<- bool) {
	for {
		if r.pgx == nil {
			// db disconnected
			if reconnectable := r.connectListener(); !reconnectable {
				// db not supported, bye
				return
			}
			if r.pgx == nil {
				// still not connected
				time.Sleep(10 * time.Second)
				continue
			}
		}

		notification, err := r.pgx.WaitForNotification(context.Background())
		if err != nil {
			log.WithError(err).Error("WaitForNotification failed, closing db connection")
			r.pgx.Close(context.Background())
			r.pgx = nil
			continue
		}

		log.WithField("payload", notification.Payload).Info("AddFeedListener got notification")

		// decode payload
		payload := postgresNotificationPayload{}
		err = json.Unmarshal([]byte(notification.Payload), &payload)
		if err != nil {
			log.WithError(err).Error("couldn't deserialize notification payload")
			continue
		}

		switch payload.Action {
		case "INSERT":
			addedFeeds <- payload.FeedID
		case "UPDATE":
			fallthrough
		case "DELETE":
			fallthrough
		case "TRUNCATE":
			needRehash <- true
		default:
			log.WithField("action", payload.Action).Warn("ignoring payload with unhandled action")
		}
	}
}
