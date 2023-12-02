package main

import (
	"fmt"
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	rabbitMqAddr := os.Getenv("RABBITMQ_ADDR")
	if rabbitMqAddr == "" {
		rabbitMqAddr = "localhost:5672"
	}

	conn, dialErr := amqp.Dial(fmt.Sprintf("amqp://guest:guest@%s", rabbitMqAddr))
	if dialErr != nil {
		log.Println(dialErr)
		return
	}
	defer conn.Close()

	log.Println("Successfully connected to RabbitMQ")

	ch, connErr := conn.Channel()
	if connErr != nil {
		log.Println(connErr)
		return
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
		log.Println(consumeErr)
		return
	}

	forever := make(chan bool)
	go func() {
		for msg := range msgs {
			log.Printf("Recieved message: %s\n", msg.Body)
		}
	}()
	<-forever
}
