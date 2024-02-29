package mq

import (
	"context"
	"filestore-server/config"
	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
	"time"
)

// Send 发送消息
func Send(config config.RabbitConf, msg []byte) error {
	// 1. 连接rabbit
	conn, err := amqp.Dial(config.RabbitURL)
	if err != nil {
		log.Error("connect rabbit failed, err:" + err.Error())
		return err
	}

	// 2. 创建一个channel
	channel, err := conn.Channel()
	if err != nil {
		log.Error("create a channel failed, err:" + err.Error())
		return err
	}

	// 3. 创建一个队列
	q, err := channel.QueueDeclare(
		config.TransOSSErrQueueName,
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		log.Error("create a queue failed, err:" + err.Error())
		return err
	}

	// 4. 发送消息
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = channel.PublishWithContext(ctx,
		config.TransExchangeName,
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        msg,
		})
	if err != nil {
		log.Error("publish a message failed, err:" + err.Error())
		return err
	}
	return nil
}

// Receive 接收消息
func Receive(config config.RabbitConf, msgChan chan []byte) error {
	// 1. 连接rabbit
	conn, err := amqp.Dial(config.RabbitURL)
	if err != nil {
		log.Error("connect rabbit failed, err:" + err.Error())
		return err
	}

	// 2. 创建一个channel
	channel, err := conn.Channel()
	if err != nil {
		log.Error("create a channel failed, err:" + err.Error())
		return err
	}

	// 3. 创建一个队列
	q, err := channel.QueueDeclare(
		config.TransOSSErrQueueName,
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		log.Error("create a queue failed, err:" + err.Error())
		return err
	}

	// 4. 创建一个消费通道
	msgs, err := channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		log.Error("create a consume failed, err:" + err.Error())
		return err
	}

	// 5. 开始循环消费消息
	for {
		select {
		case d := <-msgs:
			if len(d.Body) == 0 {
				continue
			}
			log.Printf("Received a message from %s : [%s]\n", q.Name, string(d.Body))
			msgChan <- d.Body
		}
	}

	return nil
}
