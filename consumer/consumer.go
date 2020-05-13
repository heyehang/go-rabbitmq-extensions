package consumer

import (
	"log"
	"time"

	"github.com/heyehang/go-rabbitmq-extensions/connection"

	"github.com/streadway/amqp"
)

type consumerOption struct {
	queue         string
	connectionKey connection.ConnectionKey
	consumerFunc  func(msgjson string) error
	ch            *amqp.Channel
	msgs          <-chan amqp.Delivery
}

func (opt *consumerOption) failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("queue:%s, %s: %s", opt.queue, msg, err)
	}
}

//Consume 开始消费
func (opt *consumerOption) Consume(consumerFunc func(msgjson string) error) {
	opt.consumerFunc = consumerFunc

	conn, err := opt.connectionKey.SingletonConnection()
	opt.failOnError(err, "Consume Failed to open a connection")

	ch, err := conn.Channel()
	opt.failOnError(err, "Consume Failed to open a channel")
	opt.ch = ch
	go opt.channelMonitor()
	defer ch.Close()

	q, err := ch.QueueDeclare(
		opt.queue, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	opt.failOnError(err, "Consume Failed to declare a queue")

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	opt.failOnError(err, "Consume Failed to set QoS")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	opt.msgs = msgs
	opt.failOnError(err, "Consume Failed to register a consumer")

	log.Printf(" [*] Listen queue:%s", opt.queue)

	for msg := range opt.msgs {

		err := consumerFunc(string(msg.Body))

		if err != nil {
			msg.Reject(true)
			<-time.After(1 * time.Second)
		} else {
			msg.Ack(false)
		}
	}
	log.Printf(" [*]connection or channel  is close, ReListen queue:%s", opt.queue)
}

func (opt *consumerOption) channelMonitor() {
	notifyClose := make(chan *amqp.Error)
	opt.ch.NotifyClose(notifyClose)
	select {
	case <-notifyClose:
		opt.Consume(opt.consumerFunc)
	}
}
