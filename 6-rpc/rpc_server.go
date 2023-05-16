package main

import (
	"context"
	"strconv"
	"time"

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
		"rpc_queue",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.Error().Msgf("%s: %s", "Failed to declare a queue", err.Error())
	}

	if err = ch.Qos(
		1,
		0,
		false,
	); err != nil {
		logger.Error().Msgf("%s: %s", "Failed to set QoS", err.Error())
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
		logger.Error().Msgf("%s: %s", "Failed to register a consumer")
	}

	var forever chan struct{}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		for d := range msgs {
			n, err := strconv.Atoi(string(d.Body))
			if err != nil {
				logger.Error().Msgf("%s: %s", "Failed to convert body to integer", err.Error())
			}

			logger.Info().Msgf(" [.] fib(%d)\n", n)
			response := fib(n)

			if err = ch.PublishWithContext(ctx,
				"",
				d.ReplyTo,
				false,
				false,
				amqp.Publishing{
					ContentType:   "text/plain",
					CorrelationId: d.CorrelationId,
					Body:          []byte(strconv.Itoa(response)),
				}); err != nil {
				logger.Error().Msgf("%s: %s", "Failed to publish a message")
			}

			d.Ack(false)
		}
	}()

	logger.Info().Msgf(" [*] Awaiting RPC requests\n")
	<-forever
}

func fib(n int) int {
	if n == 0 {
		return 0
	} else if n == 1 {
		return 1
	} else {
		return fib(n-1) + fib(n-2)
	}
}
