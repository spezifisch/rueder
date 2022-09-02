package controller

import (
	"github.com/gofrs/uuid"
	"github.com/spezifisch/rueder3/backend/pkg/helpers"
)

// Repository stores everything for the frontend API
type Repository interface {
	GetArticle(id uuid.UUID) (Article, error)
	GetArticles(feedID uuid.UUID, limit int, offset int) ([]ArticlePreview, error)

	Feeds() ([]Feed, error)
	GetFeed(id uuid.UUID) (Feed, error)
	GetFeedByURL(url string) (Feed, error)
	AddFeed(url string) (feedID uuid.UUID, err error)

	// tied to the user:
	Folders(*helpers.AuthClaims) ([]Folder, error)
	ChangeFolders(*helpers.AuthClaims, []Folder) error
}

// UserEventRepository can send live events to users
type UserEventRepository interface {
	Publisher() *UserEventPublisher

	HandleEvents()
	Close()
}
