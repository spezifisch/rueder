package rabbitmq

import (
	"context"

	"github.com/apex/log"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spezifisch/rueder3/backend/pkg/api/controller"
)

// EventPublisherRepository internal state
type EventPublisherRepository struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	context    context.Context

	eventInput          <-chan controller.UserEventEnvelope
	eventClose          chan struct{}
	eventHandlerRunning bool

	producer controller.UserEventPublisher
}

// NewEventPublisherRepository connects to RabbitMQ and configures an event publisher
func NewEventPublisherRepository(addr string) *EventPublisherRepository {
	ctx := context.Background()
	r, err := rabbitMQConnect(addr)
	if err != nil {
		return nil
	}

	err = r.declareUserEventsExchange()
	if err != nil {
		return nil
	}

	ownChannel := make(chan controller.UserEventEnvelope)
	closeChannel := make(chan struct{})
	return &EventPublisherRepository{
		// channels for external use by backend api
		producer: controller.UserEventPublisher{
			Channel: ownChannel,
		},

		// internal endpoints of channels
		eventInput:          ownChannel,
		eventClose:          closeChannel,
		eventHandlerRunning: false,

		connection: r.connection,
		channel:    r.channel,
		context:    ctx,
	}
}

func (r *EventPublisherRepository) Publisher() *controller.UserEventPublisher {
	return &r.producer
}

// HandleEvents should be run as a goroutine to handle passing messages to rabbitmq
func (r *EventPublisherRepository) HandleEvents() {
	if r.eventHandlerRunning {
		panic("don't run multiple instances of this function")
	}

	r.eventHandlerRunning = true
	for {
		select {
		case envelope := <-r.eventInput:
			// we receive an event from backend api, send it to rabbitmq
			err := r.publishUserEvent(envelope)
			if err != nil {
				log.WithError(err).WithField("userID", envelope.UserID).Error("couldn't publish message")
			}
		case <-r.eventClose:
			r.eventHandlerRunning = false
			return
		}
	}
}

func (r *EventPublisherRepository) publishUserEvent(envelope controller.UserEventEnvelope) (err error) {
	// use userid as routing key so all connected clients (if any) of this user receive the same event
	routingKey := envelope.UserID.String()
	// TODO define a data type for this, needs to be the same as in the frontend
	messageBody := envelope.Message

	err = r.channel.PublishWithContext(
		r.context,
		"user_events", // exchange
		routingKey,    // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        messageBody,
		})
	return
}

// Close connection to rabbitmq
func (r *EventPublisherRepository) Close() {
	// terminate HandleEvents loop
	if r.eventHandlerRunning {
		r.eventClose <- struct{}{}
	}

	if r.channel != nil {
		r.channel.Close()
		r.channel = nil
	}

	if r.connection != nil {
		r.connection.Close()
		r.connection = nil
	}
}
