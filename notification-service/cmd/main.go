package main

import (
	"log"

	"notification-service/internal/consumer"
)

func main() {

	c, err := consumer.NewConsumer(
		"amqp://guest:guest@rabbitmq:5672/",
	)
	if err != nil {
		log.Fatal(err)
	}

	err = c.Consume()
	if err != nil {
		log.Fatal(err)
	}
}
