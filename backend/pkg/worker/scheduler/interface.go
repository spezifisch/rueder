package scheduler

import (
	"github.com/gofrs/uuid"
)

// Repository is the database used by the scheduler and workers
type Repository interface {
	// -> for scheduler
	// Feeds returns all feeds that need to be fetched periodically
	Feeds() ([]Feed, error)
	GetFeed(feedID uuid.UUID) (Feed, error)
	// RunAddFeedListener starts a blocking listener that outputs at the addedFeeds channel
	// whenever a feed is added to the repository. It outputs at the needRehash channel
	// whenever a feed is removed or its parameters are changed.
	RunFeedChangeListener(addedFeeds chan<- uuid.UUID, needRehash chan<- bool) (err error)

	// -> for workers
	// UpdateFeedInfo updates the feed with the given uuid with new data from the Feed object
	UpdateFeedInfo(feedID uuid.UUID, updatedFeed *Feed) (err error)
	// CheckExistingArticles returns a list whose elements are true when the corresponding article uuid exists
	CheckExistingArticles(feedID uuid.UUID, articleGUIDs []string) (exists []bool, err error)
	// AddArticle adds a new article and associates it with the given feed
	AddArticle(feedID uuid.UUID, article *Article) error
}

// WorkerPool spawns workers that fetch feeds
type WorkerPool interface {
	// StartWorker starts a feed fetcher
	StartWorker(id int, feeds <-chan Feed, doneFeeds chan<- Feed)
}
