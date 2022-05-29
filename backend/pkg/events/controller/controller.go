package controller

// Controller for Events API v1
type Controller struct {
	repository Repository
}

// NewController for Events API v1
func NewController(repository Repository) *Controller {
	return &Controller{
		repository: repository,
	}
}
