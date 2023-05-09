package main

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal().Msgf("err: %s", err.Error())
	}
	defer conn.Close()

	fmt.Println("successfully connected to rabbitmq instance")

	ch, err := conn.Channel()
	if err != nil {
		log.Error().Msgf("err: %s", err.Error())
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"test-queue",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Error().Msgf("err: %s", err.Error())
	}

	fmt.Println(q)

	if err := ch.Publish(
		"",
		"test-queue",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte("new test message 2"),
		},
	); err != nil {
		log.Error().Msgf("err: %s", err.Error())
	}

	fmt.Println("successfully published message to rabbitmq")
}
