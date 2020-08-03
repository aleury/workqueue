package messaging

import "github.com/streadway/amqp"

type HandlerFunc func(delivery amqp.Delivery)

type Client interface {
	PublishOnQueue(msg []byte, queueName string) error
	SubscribeToQueue(queueName string, handlerFunc HandlerFunc) error
	Close()
}
