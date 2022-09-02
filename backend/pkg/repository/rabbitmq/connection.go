package rabbitmq

import (
	"github.com/apex/log"
	amqp "github.com/rabbitmq/amqp091-go"
)

// struct for setup code shared between publisher/consumer side
type rabbitMQConnection struct {
	connection *amqp.Connection
	channel    *amqp.Channel
}

// rabbitMQConnect connects and creates a single channel for events
func rabbitMQConnect(addr string) (ret rabbitMQConnection, err error) {
	log.Info("connecting to rabbitmq")

	conn, err := amqp.Dial(addr)
	if err != nil {
		log.WithError(err).WithField("addr", addr).Error("couldn't connect to rabbitmq")
		return
	}
	ret.connection = conn
	log.Info("connected to rabbitmq")

	ch, err := conn.Channel()
	if err != nil {
		log.WithError(err).WithField("addr", addr).Error("couldn't open channel")
		return
	}
	ret.channel = ch

	return
}

// declareUserEventsExchange creates the exchange between rueder api and events api
func (r *rabbitMQConnection) declareUserEventsExchange() (err error) {
	err = r.channel.ExchangeDeclare(
		"user_events", // name
		"fanout",      // type
		false,         // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		log.WithError(err).Error("couldn't declare exchange")
		return nil
	}

	return
}
