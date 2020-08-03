package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"

	"workqueue"
	"workqueue/messaging"
)

func main() {
	wait := make(chan os.Signal, 1)
	signal.Notify(wait, os.Interrupt)

	client, err := messaging.NewRabbitMQClient(workqueue.Config.AMQPConnectionURL)
	if err != nil {
		log.Fatalf("couldn't establish connection to RabbitMQ: %s", err)
	}
	defer client.Close()

	consumer := messaging.NewConsumer("add", client)
	err = consumer.OnMsg(handleAddTask)
	if err != nil {
		log.Fatalf("couldn't subscribe to the `add` queue: %s", err)
	}

	<-wait
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
