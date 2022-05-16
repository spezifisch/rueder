package models

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/gofrs/uuid"
)

// Folder is a folder with feeds belonging to a user
type Folder struct {
	ID uuid.UUID `json:"id"`

	Title string      `json:"title"`
	Feeds []LightFeed `json:"feeds"`
}

// Folders is a folder list
type Folders []Folder

// Value implements the driver.Valuer interface
func (f Folders) Value() (driver.Value, error) {
	return json.Marshal(f)
}

// Scan implements the sql.Scanner interface
func (f *Folders) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	var data = value.([]uint8)
	return json.Unmarshal(data, &f)
}

// LightFeed for the folder view
type LightFeed struct {
	ID uuid.UUID `json:"id"`

	Title string `json:"title,omitempty"`
	Icon  string `json:"icon,omitempty"`
}
