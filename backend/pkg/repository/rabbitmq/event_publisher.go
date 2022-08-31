package rabbitmq

import (
	"os"
	"strings"

	"github.com/apex/log"
	amqp "github.com/rabbitmq/amqp091-go"
)

// EventPublisherRepository internal state
type EventPublisherRepository struct {
	conn *amqp.Connection
}

// NewEventPublisherRepository connects to RabbitMQ and configures an event publisher
func NewEventPublisherRepository(addr string) *EventPublisherRepository {
	repo := EventPublisherRepository{}
	log.Info("connecting to rabbitmq")

	conn, err := amqp.Dial(addr)
	if err != nil {
		log.WithError(err).WithField("addr", addr).Error("couldn't connect to rabbitmq")
		return nil
	}
	repo.conn = conn
	log.Info("connected to rabbitmq")

	ch, err := conn.Channel()
	if err != nil {
		log.WithError(err).WithField("addr", addr).Error("couldn't open channel")
		return nil
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs_topic", // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		log.WithError(err).WithField("addr", addr).Error("couldn't declare exchange")
		return nil
	}

	body := bodyFrom(os.Args)
	err = ch.Publish(
		"logs_topic",          // exchange
		severityFrom(os.Args), // routing key
		false,                 // mandatory
		false,                 // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	if err != nil {
		log.WithError(err).WithField("addr", addr).Error("couldn't publish message")
		return nil
	}

	log.Infof("sent %s", body)

	return &repo
}

func bodyFrom(args []string) string {
	var s string
	if (len(args) < 3) || os.Args[2] == "" {
		s = "hello"
	} else {
		s = strings.Join(args[2:], " ")
	}
	return s
}

func severityFrom(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "anonymous.info"
	} else {
		s = os.Args[1]
	}
	return s
}

func (r *EventPublisherRepository) Close() {
	if r.conn != nil {
		r.conn.Close()
		r.conn = nil
	}
}
