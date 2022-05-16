package controller

// Controller for API v1
type Controller struct {
	repository      Repository
	articlesPerPage int
}

// NewController for API v1
func NewController(repository Repository) *Controller {
	return &Controller{
		repository:      repository,
		articlesPerPage: 40,
	}
}
