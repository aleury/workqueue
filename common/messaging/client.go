package messaging

type HandlerFunc func(msg []byte) error

type Client interface {
	PublishOnQueue(msg []byte, queueName string) error
	SubscribeToQueue(queueName string, handlerFunc HandlerFunc) error
	Close()
}
