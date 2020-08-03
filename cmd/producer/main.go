package main

import (
	"encoding/json"
	"flag"
	"log"
	"math/rand"
	"time"

	"workqueue"
	"workqueue/common/messaging"
)

func main() {
	ratePtr := flag.Int64("rate", 1.0, "messages per second")
	flag.Parse()

	rand.Seed(time.Now().UnixNano())

	client, err := messaging.NewRabbitMQClient(workqueue.Config.AMQPConnectionURL)
	handleErr(err, "couldn't establish connection to RabbitMQ")
	defer client.Close()

	for {
		task := mkAddTask()
		err = client.PublishOnQueue(task, "add")
		handleErr(err, "failed to publish task to the `add` queue")
		time.Sleep(time.Second / time.Duration(*ratePtr))
	}
}

func mkAddTask() []byte {
	addTask := workqueue.AddTask{Number1: rand.Intn(999), Number2: rand.Intn(999)}

	task, err := json.Marshal(addTask)
	handleErr(err, "error encoding json")

	log.Printf("Publishing %s", addTask)
	return task
}

func handleErr(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
