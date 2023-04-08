package models

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/gofrs/uuid"
)

// UserState contains feed and article state info shared between all clients of the same user
type UserState struct {
	ID int `json:"-" db:"id"`

	UserID uuid.UUID `json:"user_id" db:"user_id"`
	User   *User     `json:"user" belongs_to:"user"`

	FeedStates map[uuid.UUID]UserFeedState `json:"feed_states"`
}

// Value implements the driver.Value interface
func (u UserState) Value() (driver.Value, error) {
	return json.Marshal(u)
}

// Scan implements the sql.Scanner interface
func (u *UserState) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	var data = value.([]uint8)
	return json.Unmarshal(data, &u)
}

// UserFeedState stores all information regarding the feed that's individual to the user
type UserFeedState struct {
	FeedID uuid.UUID `json:"feed_id"`

	// the following fields all refer to article sequence numbers (the n'th article in the feed when ordered by CreatedAt)
	LastUpdateAt int   `json:"last_update_at"`
	ReadAllUntil int   `json:"read_all_until"`
	ReadArticles []int `json:"read_articles"`
}
