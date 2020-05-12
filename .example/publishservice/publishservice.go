package publishservice

import (
	"github.com/heyehang/go-rabbitmq-extensions/.example/rabbitmqoptions"

	"github.com/heyehang/go-rabbitmq-extensions/publish"
)

func APublish() *publish.PublishOption {
	return &publish.PublishOption{
		ConnectionKey: rabbitmqoptions.MyRabbitmqOptionConnKey,
		ExchangeType:  "direct",
		Exchange:      "atest.ex",
		RoutingKey:    "atest.key",
		Queue:         "atest.queue",
	}
}

func BPublish() *publish.PublishOption {
	return &publish.PublishOption{
		ConnectionKey: rabbitmqoptions.MyRabbitmqOptionConnKey,
		ExchangeType:  "direct",
		Exchange:      "btest.ex",
		RoutingKey:    "btest.key",
		Queue:         "btest.queue",
	}
}

func CPublish() *publish.PublishOption {
	return &publish.PublishOption{
		ConnectionKey: rabbitmqoptions.MyRabbitmqOptionConnKey,
		ExchangeType:  "direct",
		Exchange:      "ctest.ex",
		RoutingKey:    "ctest.key",
		Queue:         "ctest.queue",
	}
}
