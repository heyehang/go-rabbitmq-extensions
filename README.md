
# go-rabbitmq-extensions介绍

* 这是一个 基于amqp的Rabbit轻量级框架，可以让开发人员无需关注底层变动，专注编写业务代码，从而达到便捷开发。

# 特性

* go-rabbitmq-extensions，非常的精简小巧实用，下面将介绍 go-rabbitmq-extensions 的项目框架。
* 开发设计思路是将Rabbit的连接池，生产者，消费者三种业务类型分层分离，从而实现解耦轻量化。
* 连接池，生产者，消费者的设计实现逻辑采用空接口抽象,实现各自之间单一职责与开闭原则，是非常有利于业务的扩展和维护。
* 连接池：拥有网络断开的自动重连机制，并发安全的内置连接池管理。在这里用户只需要关心配置连接池相关参数。
* 生产者和消费者：底层已经全部抽象实现，无须关注底层逻辑，在这里用户只需要关心配置生产者/消费者相关参数，并且消费者支持单例多重消费者。
* 开发人员只需要在Rabbit管控台新建相关的VHost，其他参数（Exchange，Queue，ExchangeType,RoutingKey）全部代码自动帮你建立完好，无须手动新建，解决繁琐操作。
* 项目 gitbhub 地址：<https://github.com/heyehang/go-rabbitmq-extensions>

---


## 参数说明

*  HostName，Rabbit所在服务器地址
*  port,端口号
*  username，登录账号名称
*  password，登录密码
*  VHost，虚拟主机
*  Exchange，交换机
*  ExchangeType,交换机类型
*  Queue，队列名称
*  RoutingKey，队列与交换机绑定的key
*  ConnectionKey,当前连接池服务的key，推荐：(当前连接池结构体变量名称)
*  ConsumerTotal,当前消费队列所对应的消费者数量（支持单例多重消费者）
# 如何开始？


```
安装命令：go get https://github.com/heyehang/go-rabbitmq-extensions
```
    

## 连接池

* 声明类型为connection.RabbitmqOptionBase的指针变量 并确保ConnectionKey的唯一性。

## 示例代码

``` go
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

```

---

## 生产者

* 声明类型为publish.PublishOption的指针变量,并实现相关参数，并且绑定所需要使用的连接池ConnectionKey，
 发送队列消息：APublish.Publish(objmsg)

## 示例代码

``` go
package publishservice

import (
	"go-rabbitmq-extensions/rabbitmqoptions"

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
```

---

## 消费者

* 声明一个自定义结构体，并实现consumer.ConsumerService的两个方法，{ConsumerOption() (Queue string, ConnectionKey connection.ConnectionKey, ConsumerTotal int)}
和{ConsumerFunc(msgjson string) error }
## 示例代码

``` go

type AConsumeOpt struct{}

func (c *AConsumeOpt) ConsumerOption() (Queue string, ConnectionKey connection.ConnectionKey, ConsumerTotal int) {
	return "atest.queue", rabbitmqoptions.MyRabbitmqOptionConnKey, 1

}

func (c *AConsumeOpt) ConsumerFunc(msgjson string) error {
	fmt.Println(msgjson)
	return nil
}
```

---

## 服务注册


## 示例代码

``` go

func TestAll(t *testing.T) {
        //注册连接池
	rabbitmqoptions.MyRabbitmqOption.RegisterConnection()

	//生产者 发送消息
	publishservice.APublish().Publish("APublishOpt" + time.Now().String())
			
	//注册消费者
	consumerService := consumer.New()
	consumerService.RegisterConsumer(&consumerservice.AConsumeOpt{})
	consumerService.RegisterConsumer(&consumerservice.BConsumeOpt{})
	consumerService.RegisterConsumer(&consumerservice.CConsumeOpt{})
	//启动消费监听
	consumerService.Start()

	c := make(chan bool)
	<-c
}
```

---
