package main

import (
	"os"

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

	if err := ch.ExchangeDeclare(
		"logs_topic",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		logger.Error().Msgf("%s: %s", "Failed to declare an exchange", err.Error())
	}

	q, err := ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		logger.Error().Msgf("%s: %s", "Failed to declare a queue", err.Error())
	}

	if len(os.Args) < 2 {
		logger.Debug().Msgf("Usage: %s [info] [warning] [error]", os.Args[0])
		os.Exit(0)
	}

	for _, s := range os.Args[1:] {
		logger.Info().Msgf("Binding queue %s to exchange %s with routing key %s",
			q.Name, "logs_topic", s)
		if err = ch.QueueBind(
			q.Name,
			s,
			"logs_topic",
			false,
			nil,
		); err != nil {
			logger.Error().Msgf("%s: %s", "Failed to bind a queue", err.Error())
		}
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
		logger.Error().Msgf("%s: %s", "Failed to register a consumer", err)
	}

	var forever chan struct{}

	go func() {
		for d := range msgs {
			logger.Info().Msgf(" [x] %s\n", d.Body)
		}
	}()

	logger.Info().Msgf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
