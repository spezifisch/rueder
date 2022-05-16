package scheduler

import (
	"time"

	"github.com/gofrs/uuid"
)

// Article of a feed for the worker
type Article struct {
	// our own UUID for this article
	ID uuid.UUID `json:"id"`

	// the article's GUID or permalink according to the feed
	SiteGUID string `json:"site_guid"`

	Time time.Time `json:"time"`
	Link string    `json:"link,omitempty"`

	Image      string `json:"image,omitempty"`
	ImageTitle string `json:"image_title,omitempty"`

	// raw and html-sanitized versions
	Title     string `json:"title,omitempty"`
	RawTitle  string `json:"-"`
	Teaser    string `json:"teaser,omitempty"`
	RawTeaser string `json:"-"`
	Text      string `json:"text,omitempty"`
	RawText   string `json:"-"`

	Authors []string `json:"authors,omitempty"`
	Tags    []string `json:"tags,omitempty"`

	Enclosures []ArticleEnclosure `json:"enclosure,omitempty"`
}

// ArticleEnclosure is an attached file
type ArticleEnclosure struct {
	Length string `json:"length,omitempty"`
	Type   string `json:"type,omitempty"`
	URL    string `json:"url,omitempty"`
}

// FeedFetcherState contains info about a feed for the fetcher
type FeedFetcherState struct {
	// last tried fetch
	FetchedAt time.Time `json:"time,omitempty"`
	// min. delay in seconds between fetches
	FetchDelayS int `json:"fetch_delay_s,omitempty"`

	// caching headers
	ETag         string `json:"etag,omitempty"`
	LastModified string `json:"last_modified,omitempty"`

	// status
	Working     bool      `json:"working"`
	LastSuccess time.Time `json:"last_success,omitempty"`
	LastError   time.Time `json:"last_error,omitempty"`
	Message     string    `json:"message,omitempty"`
}

// Feed for the worker and scheduler
type Feed struct {
	ID           uuid.UUID `json:"id,omitempty"`
	CreatedAt    time.Time `json:"created,omitempty"`
	ArticleCount int       `json:"article_count,omitempty"`

	FetcherState FeedFetcherState `json:"fetcher_state,omitempty"`

	FeedURL string `json:"feed_url"`
	SiteURL string `json:"site_url,omitempty"`

	Title    string    `json:"title,omitempty"`
	Icon     string    `json:"icon,omitempty"`
	Articles []Article `json:"articles,omitempty"`
}
