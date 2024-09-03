package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
)

func main() {
	godotenv.Load()

	// Define RabbitMQ server URL
	amqpServerURL := os.Getenv("AMQP_SERVER_URL")

	// Create a new RabbitMQ connection
	connectRabbitMQ, err := amqp.Dial(amqpServerURL)

	if err != nil {
		panic(err)
	}
	defer connectRabbitMQ.Close()

	// Opening a channel to RabbitMQ
	// the connection we have already established
	channelRabbitMQ, err := connectRabbitMQ.Channel()
	if err != nil {
		panic(err)
	}
	defer channelRabbitMQ.Close()

	// Subcribe to the queue to consume messages
	msgs, err := channelRabbitMQ.Consume("Queue-1", "", true, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	// Build a welcome message
	welcomeMessage := "Welcome to the Go Fiber RabbitMQ Consumer"
	log.Println(welcomeMessage)
	log.Println("Waiting for messages...")

	// Make a channel to receive messages into go routine
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Printf("Received message: %s\n", d.Body)
		}
	}()

	<-forever
}
