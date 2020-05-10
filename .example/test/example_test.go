package test

import (
	"go-rabbitmq-extensions/consumer"
	"testing"
	"time"

	"go-rabbitmq-extensions/.example/consumerservice"
	"go-rabbitmq-extensions/.example/publishservice"
	"go-rabbitmq-extensions/.example/rabbitmqoptions"
)

func TestAll(t *testing.T) {
	rabbitmqoptions.MyRabbitmqOption.RegisterConnection()

	go func() {
		for {
			// go publishservice.APublishOpt.Publish("APublishOpt")
			// go publishservice.BPublishOpt.Publish("BPublishOpt")
			// go publishservice.CPublishOpt.Publish("CPublishOpt")
			//并发 for go CPU 满负荷 100% 会烧坏 采用其他方式
			publishservice.APublish().Publish("APublishOpt" + time.Now().String())
			publishservice.BPublish().Publish("BPublishOpt" + time.Now().String())
			publishservice.CPublish().Publish("CPublishOpt" + time.Now().String())
		}
	}()
	consumerService := consumer.New()
	consumerService.RegisterConsumer(&consumerservice.AConsumeOpt{})
	consumerService.RegisterConsumer(&consumerservice.BConsumeOpt{})
	consumerService.RegisterConsumer(&consumerservice.CConsumeOpt{})
	consumerService.Start()

	c := make(chan bool)
	<-c
}

func TestConsumerTotal(t *testing.T) {
	rabbitmqoptions.MyRabbitmqOption.RegisterConnection()

	go func() {
		for {
			publishservice.APublish().Publish("APublishOpt" + time.Now().String())

		}
	}()
	consumerService := consumer.New()
	consumerService.RegisterConsumer(&consumerservice.AConsumeOpt{})
	consumerService.Start()

	c := make(chan bool)
	<-c
}

func TestConnectionMonitor(t *testing.T) {
	rabbitmqoptions.MyRabbitmqOption.RegisterConnection()

	consumerService := consumer.New()
	consumerService.RegisterConsumer(&consumerservice.CConsumeOpt{})
	consumerService.Start()

	c := make(chan bool)
	<-c
}

func TestChannelMonitor(t *testing.T) {
	rabbitmqoptions.MyRabbitmqOption.RegisterConnection()

	consumerService := consumer.New()
	consumerService.RegisterConsumer(&consumerservice.CConsumeOpt{})
	consumerService.Start()

	c := make(chan bool)
	<-c
}
