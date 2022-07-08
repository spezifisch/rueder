package controller

// Controller for Events API v1
type Controller struct {
	ruederRepo RuederRepository
	eventRepo  EventRepository
}

// NewController for Events API v1
func NewController(ruederRepo RuederRepository, eventRepo EventRepository) *Controller {
	return &Controller{
		ruederRepo: ruederRepo,
		eventRepo:  eventRepo,
	}
}
