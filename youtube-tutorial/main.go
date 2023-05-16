package main

import (
	"fmt"

	log "github.com/rs/zerolog"
	"github.com/streadway/amqp"
)

func main() {
	logger := log.Logger{}

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		logger.Fatal().Msgf("err: %s", err.Error())
	}
	defer conn.Close()

	logger.Info().Msg("successfully connected to rabbitmq instance")

	ch, err := conn.Channel()
	if err != nil {
		logger.Error().Msgf("err: %s", err.Error())
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
		logger.Error().Msgf("err: %s", err.Error())
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
		logger.Error().Msgf("err: %s", err.Error())
	}

	logger.Info().Msg("successfully published message to rabbitmq")
}
