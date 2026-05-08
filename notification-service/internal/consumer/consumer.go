package consumer

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn *amqp.Connection
	ch   *amqp.Channel

	processed map[string]bool
}

func NewConsumer(url string) (*Consumer, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	_, err = ch.QueueDeclare(
		"payment.completed",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		conn:      conn,
		ch:        ch,
		processed: make(map[string]bool),
	}, nil
}

func (c *Consumer) Consume() error {

	msgs, err := c.ch.Consume(
		"payment.completed",
		"",
		false, // autoAck = false
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	log.Println("	Waiting for messages...")

	for msg := range msgs {

		body := string(msg.Body)

		// idempotency check
		if c.processed[body] {
			log.Println("⚠Duplicate message skipped:", body)

			msg.Ack(false)
			continue
		}

		// mark as processed
		c.processed[body] = true

		log.Println("Notification received:", body)

		// manual ACK
		msg.Ack(false)

		log.Println("Message ACKed")
	}
	return nil
}
