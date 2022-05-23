package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"
)

// Article is used by pop to map your articles database table to your go code.
type Article struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Seq       int       `json:"seq" db:"seq" rw:"r"`           // this is a BIGSERIAL in postgres (auto incrementing)
	FeedSeq   int       `json:"feed_seq" db:"feed_seq" rw:"r"` // sequence number of article per feed

	FeedID uuid.UUID `json:"-" db:"feed_id"`
	Feed   *Feed     `json:"feed" belongs_to:"feed"`

	SiteGUID string `json:"site_guid" db:"site_guid"` // according to the RSS feed

	PostedAt time.Time `json:"posted_at" db:"posted_at"` // according to the RSS feed

	Link nulls.String `json:"link" db:"link"`

	Thumbnail  nulls.String `json:"thumbnail" db:"thumbnail"`
	Image      nulls.String `json:"image" db:"image"`
	ImageTitle nulls.String `json:"image_title" db:"image_title"`

	Title   nulls.String   `json:"title" db:"title"`
	Teaser  nulls.String   `json:"teaser" db:"teaser"`
	Content ArticleContent `json:"content" db:"content"`
}

// String is not required by pop and may be deleted
func (a Article) String() string {
	ja, _ := json.Marshal(a)
	return string(ja)
}

// Articles is not required by pop and may be deleted
type Articles []Article

// String is not required by pop and may be deleted
func (a Articles) String() string {
	ja, _ := json.Marshal(a)
	return string(ja)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (a *Article) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (a *Article) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (a *Article) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
