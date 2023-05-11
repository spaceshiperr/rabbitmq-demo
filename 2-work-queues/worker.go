package main

import (
	"bytes"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func returnErr(err error, msg string) {
	log.Panicf("%s: %s", msg, err)
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

	if err := ch.Qos(1, 0, false); err != nil {
		returnErr(err, "Failed to set QoS")
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		returnErr(err, "Failed to register a consumer")
	}

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)

			dotCount := bytes.Count(d.Body, []byte("."))

			t := time.Duration(dotCount)
			time.Sleep(t * time.Second)

			log.Printf("Done")

			if err := d.Ack(false); err != nil {
				returnErr(err, "Failed to acknowledge message consuming")
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}