package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"
)

// Feed is used by pop to map your feeds database table to your go code.
type Feed struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`

	// fields for feed fetcher
	FetchedAt   time.Time `json:"fetched_at" db:"fetched_at"`
	FetchDelayS int       `json:"fetch_delay_s" db:"fetch_delay_s"`
	// for both fetcher and frontend
	FetcherState FetcherState `json:"fetcher_state" db:"fetcher_state"`

	FeedURL string       `json:"feed_url" db:"feed_url"`
	SiteURL nulls.String `json:"site_url" db:"site_url"`

	Title nulls.String `json:"title" db:"title"`
	Icon  nulls.String `json:"icon" db:"icon"`

	Articles Articles `json:"articles" has_many:"articles" order_by:"articles.posted_at desc"`

	// users who subscribed to this feed
	Users Users `json:"users" many_to_many:"user_feeds"`
}

// String is not required by pop and may be deleted
func (f Feed) String() string {
	jf, _ := json.Marshal(f)
	return string(jf)
}

// Feeds is not required by pop and may be deleted
type Feeds []Feed

// String is not required by pop and may be deleted
func (f Feeds) String() string {
	jf, _ := json.Marshal(f)
	return string(jf)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (f *Feed) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (f *Feed) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (f *Feed) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
