package consumerservice

import (
	"fmt"
	"time"

	"github.com/heyehang/go-rabbitmq-extensions/.example/rabbitmqoptions"

	"github.com/heyehang/go-rabbitmq-extensions/connection"
)

type AConsumeOpt struct{}

func (c *AConsumeOpt) ConsumerOption() (Queue string, ConnectionKey connection.ConnectionKey, ConsumerTotal int) {
	return "atest.queue", rabbitmqoptions.MyRabbitmqOptionConnKey, 10

}

func (c *AConsumeOpt) ConsumerFunc(msgjson string) error {
	fmt.Println(msgjson)
	return nil
}

type BConsumeOpt struct{}

func (c *BConsumeOpt) ConsumerOption() (Queue string, ConnectionKey connection.ConnectionKey, ConsumerTotal int) {
	return "btest.queue", rabbitmqoptions.MyRabbitmqOptionConnKey, 1

}

func (c *BConsumeOpt) ConsumerFunc(msgjson string) error {
	fmt.Println(msgjson)
	return nil
}

type CConsumeOpt struct{}

func (c *CConsumeOpt) ConsumerOption() (Queue string, ConnectionKey connection.ConnectionKey, ConsumerTotal int) {
	return "ctest.queue", rabbitmqoptions.MyRabbitmqOptionConnKey, 1

}

func (c *CConsumeOpt) ConsumerFunc(msgjson string) error {
	fmt.Println(msgjson)
	<-time.After(3 * time.Second)
	return nil
}
