package models

import (
	"encoding/json"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"
)

// UserFeed is used by pop to map your user_feeds database table to your go code.
type UserFeed struct {
	UserID uuid.UUID `json:"-" db:"user_id"`
	User   *User     `json:"user" belongs_to:"user"`

	FeedID uuid.UUID `json:"-" db:"feed_id"`
	Feed   *Feed     `json:"feed" belongs_to:"feed"`
}

// Table gives pop the name of the database table
func (u UserFeed) Table() string {
	return "user_feeds"
}

// String is not required by pop and may be deleted
func (u UserFeed) String() string {
	ja, _ := json.Marshal(u)
	return string(ja)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (u *UserFeed) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (u *UserFeed) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (u *UserFeed) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
