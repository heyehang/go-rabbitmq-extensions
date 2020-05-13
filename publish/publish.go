package publish

import (
	"encoding/json"
	"log"

	"github.com/heyehang/go-rabbitmq-extensions/connection"

	"github.com/streadway/amqp"
)

type PublishOption struct {
	ExchangeType,
	Exchange,
	RoutingKey,
	Queue string
	ConnectionKey connection.ConnectionKey
}

func (opt *PublishOption) failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("queue:%s, %s: %s", opt.Queue, msg, err)
	}
}
func (opt *PublishOption) Publish(objectmsg interface{}) error {
	conn, err := opt.ConnectionKey.SingletonConnection()
	opt.failOnError(err, "Publish Failed to open a connection")

	ch, err := conn.Channel()
	opt.failOnError(err, "Publish Failed to open a channel")

	defer ch.Close()

	err = ch.ExchangeDeclare(opt.Exchange, opt.ExchangeType, true, false, false, false, nil)

	opt.failOnError(err, "Publish Failed to Exchange Declare a channel")

	_, err = ch.QueueDeclare(
		opt.Queue, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	opt.failOnError(err, "Publish Failed to declare a queue")

	err = ch.QueueBind(opt.Queue, opt.RoutingKey, opt.Exchange, false, nil)
	opt.failOnError(err, "Publish Failed to Exchange QueueBind a channel")

	jsonstr, err := json.Marshal(objectmsg)

	if err != nil {
		opt.failOnError(err, "Publish Failed to json.Marshal")
		return err
	}
	err = ch.Publish(
		opt.Exchange,   // exchange
		opt.RoutingKey, // routing key
		false,          // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(jsonstr),
		})

	opt.failOnError(err, "Publish Failed to publish a message")

	return nil
}
