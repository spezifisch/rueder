package controller

// Controller for API v1
type Controller struct {
	repository          Repository
	userEventRepository UserEventRepository
	articlesPerPage     int
}

// NewController for API v1
func NewController(repository Repository, userEventRepository UserEventRepository) *Controller {
	return &Controller{
		repository:          repository,
		userEventRepository: userEventRepository,
		articlesPerPage:     40,
	}
}
