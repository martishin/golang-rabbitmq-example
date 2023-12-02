package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	timeout        = 5 * time.Second
	sleepOnFailure = 5 * time.Second
	message        = "Hello World"
)

func produce() error {
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

	q, queueErr := ch.QueueDeclare(
		"TestQueue",
		false,
		false,
		false,
		false,
		nil,
	)
	log.Println(q)

	if queueErr != nil {
		return queueErr
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	id := 0
	for range ticker.C {
		msg := fmt.Sprintf("%s: %d", message, id)
		publishErr := ch.PublishWithContext(
			ctx,
			"",
			"TestQueue",
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(msg),
			},
		)

		if publishErr != nil {
			return publishErr
		}
		log.Printf("Published message: %s\n", msg)
		id++
	}

	return nil
}

func main() {
	err := produce()
	if err != nil {
		log.Fatal(err)
	}
}
