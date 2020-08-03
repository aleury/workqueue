# workqueue

Task queue application for computing the sum of two integers.

## Start RabbitMQ
```bash
$ rabbitmq-server
```

## Start a consumer
```bash
$ cd cmd/consumer
$ go build
$ ./consumer
```

## Start a producer
```bash
$ cd cmd/producer
$ go build
$ ./producer -rate 10  # i.e 10 msg/s
```