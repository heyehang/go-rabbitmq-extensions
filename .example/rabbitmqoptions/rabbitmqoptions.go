package rabbitmqoptions

import "github.com/heyehang/go-rabbitmq-extensions/connection"

const (
	MyRabbitmqOptionConnKey = "MyRabbitmqOption"
)

var MyRabbitmqOption = &connection.RabbitmqOptionBase{
	Host:          "localhost",
	Port:          "5672",
	VHost:         "test.host",
	UserName:      "guest",
	PassWord:      "guest",
	ConnectionKey: MyRabbitmqOptionConnKey,
}
