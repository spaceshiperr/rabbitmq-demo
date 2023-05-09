package main

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
)

func main() {
	fmt.Println("consumer app")

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Error().Msg(err.Error())
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Error().Msg(err.Error())
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
		log.Error().Msg(err.Error())
	}

	forever := make(chan bool)
	go func() {
		for msg := range msgs {
			fmt.Printf("reciever message: %s\n", msg.Body)
		}
	}()

	fmt.Println("successfully connected to rabbitmq instance")
	fmt.Println("waiting for messages")
	<-forever
}
