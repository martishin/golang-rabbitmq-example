package main

import (
	"fmt"
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	sleepOnFailure = 5 * time.Second
)

func consume() error {
	rabbitMqAddr := os.Getenv("RABBITMQ_ADDR")
	if rabbitMqAddr == "" {
		rabbitMqAddr = "localhost:5672"
	}

	conn, dialErr := amqp.Dial(fmt.Sprintf("amqp://guest:guest@%s", rabbitMqAddr))
	if dialErr != nil {
		time.Sleep(sleepOnFailure)
		return dialErr
	}
	defer conn.Close()

	log.Println("Successfully connected to RabbitMQ")

	ch, connErr := conn.Channel()
	if connErr != nil {
		return connErr
	}
	defer ch.Close()

	msgs, consumeErr := ch.Consume(
		"TestQueue",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if consumeErr != nil {
		return consumeErr
	}

	forever := make(chan bool)
	go func() {
		for msg := range msgs {
			log.Printf("Recieved message: %s\n", msg.Body)
		}
	}()
	<-forever

	return nil
}

func main() {
	err := consume()
	if err != nil {
		log.Fatal(err)
	}
}
