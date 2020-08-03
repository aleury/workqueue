package main

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"

	"workqueue"
	"workqueue/common/messaging"
)

func main() {
	client, err := messaging.NewRabbitMQClient(workqueue.Config.AMQPConnectionURL)
	handleErr(err, "couldn't establish connection to RabbitMQ")
	defer client.Close()

	stopChan := make(chan bool)

	consumer := messaging.NewConsumer("add", client)
	err = consumer.OnMsg(handleAddTask)
	handleErr(err, "couldn't subscribe to the `add` queue")

	<-stopChan
}

func handleAddTask(msg amqp.Delivery) {
	log.Printf("Received a message: %s\n", msg.Body)

	addTask := workqueue.AddTask{}
	err := json.Unmarshal(msg.Body, &addTask)
	if err != nil {
		log.Printf("Error decoding message json: %s\n", err)
	} else {
		log.Printf("Result %s: %d\n", addTask, addTask.Run())
	}

	if err := msg.Ack(false); err != nil {
		log.Printf("Error acknowledging message: %s\n", err)
	} else {
		log.Println("Acknowledged message.")
	}
}

func handleErr(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
