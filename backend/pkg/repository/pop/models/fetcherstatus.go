package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// FetcherState contains status info that interests both backend workers and frontend
type FetcherState struct {
	// internal
	ETag         string `json:"etag,omitempty"`
	LastModified string `json:"last_modified,omitempty"`

	// things the frontend is interested in
	Working     bool      `json:"working"`
	LastSuccess time.Time `json:"last_success,omitempty"`
	LastError   time.Time `json:"last_error,omitempty"`
	Message     string    `json:"message,omitempty"`
}

// Value implements the driver.Valuer interface
func (f FetcherState) Value() (driver.Value, error) {
	return json.Marshal(f)
}

// Scan implements the sql.Scanner interface
func (f *FetcherState) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	var data = value.([]uint8)
	return json.Unmarshal(data, &f)
}
