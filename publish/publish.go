package publish

import (
	"encoding/json"
	"log"

	"github.com/heyehang/go-rabbitmq-extensions/connection"

	"github.com/streadway/amqp"
)

func (class *PublishOption) failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

type PublishOption struct {
	ExchangeType,
	Exchange,
	RoutingKey,
	Queue string
	ConnectionKey connection.ConnectionKey
}

func (opt *PublishOption) Publish(objectmsg interface{}) error {
	conn, err := opt.ConnectionKey.SingletonConnection()
	if err != nil {
		log.Fatalf("%s: %s", "Failed to open a connection", err)
		return err
	}
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("%s: %s", "Failed to open a channel", err)
		return err
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(opt.Exchange, opt.ExchangeType, true, false, false, false, nil)

	if err != nil {
		opt.failOnError(err, "Failed to Exchange Declare a channel")
		return err
	}
	_, err = ch.QueueDeclare(
		opt.Queue, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	if err != nil {
		log.Fatalf("%s: %s", "Failed to declare a queue", err)
		return err
	}
	err = ch.QueueBind(opt.Queue, opt.RoutingKey, opt.Exchange, false, nil)
	if err != nil {
		log.Fatalf("%s: %s", "Failed to declare a queue", err)
		return err
	}
	jsonstr, err := json.Marshal(objectmsg)

	if err != nil {
		log.Fatalf("%s: %s", "Failed to publish a message", err)
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

	if err != nil {
		log.Fatalf("%s: %s", "Failed to publish a message", err)
		return err
	}

	return nil
}
