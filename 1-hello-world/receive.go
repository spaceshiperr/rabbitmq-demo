package main

import (
	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/rs/zerolog"
)

func main() {
	logger := log.Logger{}

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		logger.Fatal().Msgf("%s: %s", "Failed to connect to RabbitMQ", err.Error())
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		logger.Error().Msgf("%s: %s", "Failed to open a channel", err.Error())
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.Error().Msgf("%s: %s", "Failed to declare a queue", err.Error())
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.Error().Msgf("%s: %s", "Failed to register a consumer", err.Error())
	}

	var forever chan struct{}

	go func() {
		for d := range msgs {
			logger.Info().Msgf("Received a message: %s\n", d.Body)
		}
	}()

	logger.Info().Msgf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
