package controller

import (
	"time"

	"github.com/gofrs/uuid"
)

// Article of a feed for the content view
type Article struct {
	ID uuid.UUID `json:"id"`

	FeedID    uuid.UUID `json:"feed_id"`
	FeedTitle string    `json:"feed_title,omitempty"`
	FeedURL   string    `json:"feed_url,omitempty"`

	Title        string    `json:"title"`
	Time         time.Time `json:"time"`
	Link         string    `json:"link,omitempty"`
	LinkComments string    `json:"link_comments,omitempty"`

	Thumbnail  string `json:"thumbnail,omitempty"`
	Image      string `json:"image,omitempty"`
	ImageTitle string `json:"image_title,omitempty"`

	Teaser  string         `json:"teaser,omitempty"`
	Content ArticleContent `json:"content,omitempty"`
}

// ArticleContent is everything that's only visible when the article is open
type ArticleContent struct {
	Authors    []string           `json:"authors,omitempty"`
	Tags       []string           `json:"tags,omitempty"`
	Enclosures []ArticleEnclosure `json:"enclosures,omitempty"`
	Text       string             `json:"text,omitempty"`
}

// ArticleEnclosure is an attached file
type ArticleEnclosure struct {
	URL    string `json:"url,omitempty"`
	Length string `json:"length,omitempty"`
	Type   string `json:"type,omitempty"`
}

// ArticlePreview is the short version on an article for the feed view
type ArticlePreview struct {
	ID      uuid.UUID `json:"id"`
	Seq     int       `json:"seq"`
	FeedSeq int       `json:"feed_seq"`

	Title     string    `json:"title,omitempty"`
	Time      time.Time `json:"time"`
	FeedTitle string    `json:"feed_title,omitempty"`
	FeedIcon  string    `json:"feed_icon,omitempty"`
	Teaser    string    `json:"teaser,omitempty"`
}

// Feed for the folder view
type Feed struct {
	ID uuid.UUID `json:"id"`

	Title        string `json:"title,omitempty"`
	Icon         string `json:"icon,omitempty"`
	URL          string `json:"url,omitempty"`
	SiteURL      string `json:"site_url,omitempty"`
	ArticleCount int    `json:"article_count,omitempty"`

	FetcherState *FetcherState `json:"fetcher_state,omitempty"`

	Articles []ArticlePreview `json:"articles,omitempty"`
}

// FetcherState for feed detail view
type FetcherState struct {
	Working     bool      `json:"working"`
	LastSuccess time.Time `json:"last_success,omitempty"`
	LastError   time.Time `json:"last_error,omitempty"`
	Message     string    `json:"message,omitempty"`
}

// Folder contains feeds
type Folder struct {
	ID uuid.UUID `json:"id"`

	Title string `json:"title,omitempty"`
	Feeds []Feed `json:"feeds,omitempty"`
}
