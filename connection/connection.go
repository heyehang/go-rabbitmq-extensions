package connection

import (
	"errors"
	"sync"
	"time"

	"github.com/streadway/amqp"
)

var connectionPool = make(map[ConnectionKey]*RabbitmqOptionBase)

type ConnectionKey string

type RabbitmqOptionBase struct {
	Host,
	Port,
	VHost,
	UserName,
	PassWord string
	ConnectionKey ConnectionKey
	connection    *amqp.Connection
	mutex         sync.Mutex
}

func (opt *RabbitmqOptionBase) getConnection() (connection *amqp.Connection, err error) {
	if opt.connection != nil && !opt.connection.IsClosed() {
		return opt.connection, nil
	}
	opt.mutex.Lock()
	defer opt.mutex.Unlock()
	if opt.connection != nil && !opt.connection.IsClosed() {
		return opt.connection, nil
	}
	for {
		connection, err = amqp.DialConfig("amqp://@"+opt.Host+":"+opt.Port+"/",
			amqp.Config{
				Vhost: opt.VHost,
				SASL: []amqp.Authentication{
					&amqp.PlainAuth{
						Username: opt.UserName,
						Password: opt.PassWord,
					},
				},
			})
		if err == nil {
			break
		}
		<-time.After(1 * time.Second)
	}
	opt.connection = connection
	go opt.connectionMonitor()
	return opt.connection, err
}

func (opt *RabbitmqOptionBase) RegisterConnection() {
	connectionPool[opt.ConnectionKey] = opt
}

func (connkey ConnectionKey) SingletonConnection() (conn *amqp.Connection, err error) {
	opt, ok := connectionPool[connkey]
	if !ok {
		err = errors.New("ConnectionPool未找到相应的Connection,请确保ConnectionKey是否匹配实现")
		return
	}
	conn, err = opt.getConnection()
	return
}

func (opt *RabbitmqOptionBase) connectionMonitor() {
	notifyClose := make(chan *amqp.Error)
	opt.connection.NotifyClose(notifyClose)
	select {
	case <-notifyClose:
		opt.getConnection()
	}
}
