package controller

// Controller for Events API v1
type Controller struct {
	ruederRepo RuederRepository
	eventRepo  UserEventRepository
}

// NewController for Events API v1
func NewController(ruederRepo RuederRepository, eventRepo UserEventRepository) *Controller {
	return &Controller{
		ruederRepo: ruederRepo,
		eventRepo:  eventRepo,
	}
}
