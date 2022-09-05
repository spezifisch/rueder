package common

// UserEventMessage is the message type passed from backend api to event api using EventRepository
type UserEventMessage struct {
	Type string            `json:"type"`
	Data map[string]string `json:"data,omitempty"`
}
