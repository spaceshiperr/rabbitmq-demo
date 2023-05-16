package main

import (
	log "github.com/rs/zerolog"
	"github.com/streadway/amqp"
)

func main() {
	logger := log.Logger{}

	logger.Info().Msgf("consumer app\n")

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		logger.Error().Msg(err.Error())
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		logger.Error().Msg(err.Error())
	}
	defer ch.Close()

	msgs, err := ch.Consume(
		"test-queue",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.Error().Msg(err.Error())
	}

	forever := make(chan bool)
	go func() {
		for msg := range msgs {
			logger.Info().Msgf("receiver message: %s\n", msg.Body)
		}
	}()

	logger.Info().Msg("successfully connected to rabbitmq instance")
	logger.Info().Msg("waiting for messages")

	<-forever
}
