package consumer

import (
	"github.com/heyehang/go-rabbitmq-extensions/connection"
)

type Service struct {
	ConsumerServiceList []ConsumerService
}

func New() *Service { return &Service{} }

type ConsumerService interface {
	ConsumerOption() (Queue string, ConnectionKey connection.ConnectionKey, ConsumerTotal int)
	ConsumerFunc(msgjson string) error
}

func (list *Service) RegisterConsumer(c ConsumerService) {
	list.ConsumerServiceList = append(list.ConsumerServiceList, c)
}
func (s *Service) run(c ConsumerService) {
	queue, connkey, consumerTotal := c.ConsumerOption()
	if consumerTotal <= 0 {
		consumerTotal = 1
	}
	for i := 0; i < consumerTotal; i++ {
		go (&consumerOption{
			queue:         queue,
			connectionKey: connkey,
		}).Consume(c.ConsumerFunc)
	}

}

func (s *Service) Start() {
	for _, c := range s.ConsumerServiceList {
		go s.run(c)
	}
}
