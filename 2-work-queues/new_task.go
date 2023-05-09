package main

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func returnErr(err error, msg string) {
	log.Panicf("%s: %s", msg, err)
}

func bodyFrom(args []string) string {
	var s string

	if (len(args) < 2) || os.Args[1] == "" {
		s = "hello"
	} else {
		s = strings.Join(args[1:], " ")
	}

	return s
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		returnErr(err, "Failed to connect to RabbitMQ")
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		returnErr(err, "Failed to open a channel")
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"task_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		returnErr(err, "Failed to declare a queue")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := bodyFrom(os.Args)

	if err := ch.PublishWithContext(
		ctx,
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         []byte(body),
			DeliveryMode: amqp.Persistent,
		}); err != nil {
		returnErr(err, "Failed to publish a message")
	}

	log.Printf(" [x] Sent %s\n", body)
}
