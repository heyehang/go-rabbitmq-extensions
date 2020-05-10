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

	consumerFunc func(msgjson string) error
	ch           *amqp.Channel
	msgs         <-chan amqp.Delivery
}

func (opt *consumerOption) failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

//Consume 开始消费
func (opt *consumerOption) Consume(consumerFunc func(msgjson string) error) {
	opt.consumerFunc = consumerFunc
	conn, err := opt.connectionKey.SingletonConnection()
	opt.failOnError(err, "Failed to open a connection")
	ch, err := conn.Channel()
	opt.failOnError(err, "Failed to open a channel")

	opt.ch = ch
	go opt.ChannelMonitor()
	defer ch.Close()

	q, err := ch.QueueDeclare(
		opt.queue, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	opt.failOnError(err, "Failed to declare a queue")

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	opt.failOnError(err, "Failed to set QoS")

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
	opt.failOnError(err, "Failed to register a consumer")

	go func() {
		for msg := range opt.msgs {

			err := consumerFunc(string(msg.Body))

			if err != nil {
				msg.Reject(true)
				time.Sleep(1 * time.Second) //休息一秒
			} else {
				msg.Ack(false)
			}
		}
	}()

	log.Printf(" [*] Listen queue:%s", opt.queue)
	forever := make(chan bool)
	<-forever
}

func (opt *consumerOption) ChannelMonitor() {
	notifyClose := make(chan *amqp.Error)
	opt.ch.NotifyClose(notifyClose)
	select {
	case <-notifyClose:
		opt.Consume(opt.consumerFunc)
	}
}
