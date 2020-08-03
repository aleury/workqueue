package messaging

import (
	"errors"
	"log"

	"github.com/streadway/amqp"
)

type RabbitMQClient struct {
	conn *amqp.Connection
}

func NewRabbitMQClient(connectionString string) (*RabbitMQClient, error) {
	if connectionString == "" {
		return nil, ErrEmptyConnString
	}

	conn, err := amqp.Dial(connectionString)
	if err != nil {
		return nil, err
	}

	return &RabbitMQClient{conn: conn}, nil
}

func (c *RabbitMQClient) PublishOnQueue(msg []byte, queueName string) error {
	if c.conn == nil {
		return ErrUninitializedConn
	}

	// Get a channel from the connection.
	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}

	// Declare the queue. It will be created if it doesn't exist.
	queue, err := ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return err
	}

	// Publish a message onto the queue.
	err = ch.Publish("", queue.Name, false, false, amqp.Publishing{
		Body:         msg,
		DeliveryMode: amqp.Persistent,
		ContentType:  "application/json",
	})
	if err != nil {
		return err
	}

	// Cleanup AMQP Channel.
	err = ch.Close()
	if err != nil {
		return err
	}

	return nil
}

func (c *RabbitMQClient) SubscribeToQueue(queueName string, handlerFunc HandlerFunc) error {
	if c.conn == nil {
		return ErrUninitializedConn
	}

	// Get a channel from the connection.
	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}

	err = ch.Qos(1, 0, false)
	if err != nil {
		return err
	}

	// Declare the queue. It will be created if it doesn't exist.
	queue, err := ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	go consumeLoop(msgs, handlerFunc)

	return nil
}

func (c *RabbitMQClient) Close() {
	if c.conn == nil {
		return
	}

	err := c.conn.Close()
	if err != nil {
		log.Println(err)
	}
	c.conn = nil
}

func consumeLoop(deliveries <-chan amqp.Delivery, handlerFunc HandlerFunc) {
	for d := range deliveries {
		err := handlerFunc(d.Body)
		if err != nil {
			log.Printf("Error handling message: %s\n", err)
			err = d.Nack(false, true)
			if err != nil {
				log.Printf("Error nacking message: %s", err)
			}
			continue
		}

		if err := d.Ack(false); err != nil {
			log.Printf("Error acknowledging message: %s\n", err)
		} else {
			log.Println("Acknowledged message.")
		}
	}
}

var (
	ErrEmptyConnString   = errors.New("cannot initialize connection to broker with an empty connection string")
	ErrUninitializedConn = errors.New("attempted to use a connection that hasn't been initialized")
)
