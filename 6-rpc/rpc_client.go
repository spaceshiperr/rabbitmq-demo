package main

import (
	"context"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/rs/zerolog"
)

func main() {
	logger := log.Logger{}

	rand.Seed(time.Now().UTC().UnixNano())

	n := bodyFrom(logger, os.Args)

	logger.Info().Msgf(" [x] Requesting fib(%d)\n", n)

	res, err := fibonacciRPC(logger, n)
	if err != nil {
		logger.Error().Msgf("%s: %s", "Failed to handle RPC request", err.Error())
	}

	logger.Info().Msgf(" [.] Got %d\n", res)
}

func bodyFrom(logger log.Logger, args []string) int {
	var s string

	if (len(args) < 2) || os.Args[1] == "" {
		s = "30"
	} else {
		s = strings.Join(args[1:], " ")
	}

	n, err := strconv.Atoi(s)
	if err != nil {
		logger.Error().Msgf("%s: %s", "Failed to convert arg to integer", err.Error())
	}

	return n
}

func randomString(l int) string {
	bytes := make([]byte, l)

	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func fibonacciRPC(logger log.Logger, n int) (res int, err error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		logger.Fatal().Msgf("%s: %s", "Failed to connect to RabbitMQ")
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		logger.Error().Msgf("%s: %s", "Failed to open a channel", err.Error())
	}
	defer ch.Close()

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

	corrId := randomString(32)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = ch.PublishWithContext(ctx,
		"",
		"rpc_queue",
		false,
		false,
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrId,
			ReplyTo:       q.Name,
			Body:          []byte(strconv.Itoa(n)),
		}); err != nil {
		logger.Error().Msgf("%s: %s", "Failed to publish a message", err.Error())
	}

	for d := range msgs {
		if corrId == d.CorrelationId {
			res, err = strconv.Atoi(string(d.Body))
			if err != nil {
				logger.Error().Msgf("%s: %s", "Failed to convert body to integer", err.Error())
				break
			}
		}
	}

	return
}
