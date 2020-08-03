package main

import (
	"encoding/json"
	"log"

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

func handleAddTask(msg []byte) error {
	log.Printf("Received a message: %s\n", msg)

	addTask := workqueue.AddTask{}
	err := json.Unmarshal(msg, &addTask)
	if err != nil {
		log.Printf("Error decoding message json: %s\n", err)
		return err
	}

	log.Printf("Result %s: %d\n", addTask, addTask.Run())
	return nil
}

func handleErr(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
