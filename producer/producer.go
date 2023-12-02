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
	timeoutInSeconds = 5
	message          = "Hello World"
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
		log.Println(queueErr)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutInSeconds*time.Second)
	defer cancel()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		publishErr := ch.PublishWithContext(
			ctx,
			"",
			"TestQueue",
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(message),
			},
		)

		if publishErr != nil {
			log.Println(publishErr)
			return
		}
		log.Printf("Successfully published message: %s\n", message)
	}
}
