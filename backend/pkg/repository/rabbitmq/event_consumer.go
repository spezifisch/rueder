package rabbitmq

import (
	"github.com/apex/log"
	"github.com/gofrs/uuid"
	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/spezifisch/rueder3/backend/pkg/events/controller"
)

// EventConsumerRepository internal state
type EventConsumerRepository struct {
	connection *amqp.Connection
	channel    *amqp.Channel
}

// NewEventConsumerRepository returns a RedisRepository that wraps a redis DB
func NewEventConsumerRepository(addr string) *EventConsumerRepository {
	repo := EventConsumerRepository{}
	log.Info("connecting to rabbitmq")

	conn, err := amqp.Dial(addr)
	if err != nil {
		log.WithError(err).WithField("addr", addr).Error("couldn't connect to rabbitmq")
		return nil
	}
	repo.connection = conn
	log.Info("connected to rabbitmq")

	ch, err := conn.Channel()
	if err != nil {
		log.WithError(err).WithField("addr", addr).Error("couldn't open channel")
		return nil
	}
	repo.channel = ch

	err = ch.ExchangeDeclare(
		"user_events", // name
		"fanout",      // type
		false,         // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		log.WithError(err).WithField("addr", addr).Error("couldn't declare exchange")
		return nil
	}

	return &repo
}

func (r *EventConsumerRepository) ConnectUser(uuid uuid.UUID) (state controller.EventUserState, err error) {
	q, err := r.channel.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return
	}

	err = r.channel.QueueBind(
		q.Name,        // queue name
		uuid.String(), // routing key
		"user_events", // exchange
		false,
		nil,
	)
	if err != nil {
		return
	}

	msgs, err := r.channel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return
	}

	// wrap messages in our own data type
	ownChannel := make(chan controller.EventMessage)
	closeChannel := make(chan struct{}, 1)
	go func(queueName, userID string) {
		for {
			select {
			case d := <-msgs:
				ownChannel <- controller.EventMessage{
					Message: d.Body, // TODO does this need a deepcopy?
				}
			case <-closeChannel:
				_, err := r.channel.QueueDelete(queueName, false, false, false)
				if err != nil {
					log.WithError(err).WithField("user", userID).Error("failed deleting queue")
				}
				return
			}
		}
	}(q.Name, uuid.String())

	state = controller.EventUserState{
		Channel: ownChannel,
		Close:   closeChannel,
	}
	return
}

func (r *EventConsumerRepository) Close() {
	if r.channel != nil {
		r.channel.Close()
		r.channel = nil
	}

	if r.connection != nil {
		r.connection.Close()
		r.connection = nil
	}
}
