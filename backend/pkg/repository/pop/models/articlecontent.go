package models

import (
	"database/sql/driver"
	"encoding/json"
)

// ArticleContent is everything that's only visible when the article is open
type ArticleContent struct {
	Authors    []string           `json:"authors,omitempty"`
	Tags       []string           `json:"tags,omitempty"`
	Enclosures []ArticleEnclosure `json:"enclosures,omitempty"`
	Text       string             `json:"text,omitempty"`
}

// Value implements the driver.Valuer interface
func (a ArticleContent) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Scan implements the sql.Scanner interface
func (a *ArticleContent) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	var data = value.([]uint8)
	return json.Unmarshal(data, &a)
}

// ArticleEnclosure is an attached file
type ArticleEnclosure struct {
	URL    string `json:"url,omitempty"`
	Length string `json:"length,omitempty"`
	Type   string `json:"type,omitempty"`
}
