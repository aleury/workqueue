package messaging

type Consumer struct {
	queue  string
	client Client
}

func NewConsumer(queue string, client Client) *Consumer {
	return &Consumer{
		queue:  queue,
		client: client,
	}
}

func (c *Consumer) OnMsg(handlerFunc HandlerFunc) error {
	return c.client.SubscribeToQueue(c.queue, handlerFunc)
}
