package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
)

func main() {
	godotenv.Load()

	// Define the RabbitMQ server URL
	amqpServerUrl := os.Getenv("AMQP_SERVER_URL")

	// Create a new RabbitMQ connection
	connectRabbitMQ, err := amqp.Dial(amqpServerUrl)

	if err != nil {
		panic(err)
	}
	defer connectRabbitMQ.Close()

	// Let's start by opening a channel to our RabbitMQ
	// instance over the connection we have already
	// established.
	channelRabbitMQ, err := connectRabbitMQ.Channel()
	if err != nil {
		panic(err)
	}
	defer channelRabbitMQ.Close()

	// With the instance and declare Queues that we can
	// publish and subscribe to.
	_, err = channelRabbitMQ.QueueDeclare("Queue-1", true, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	// Create a fiber instance
	app := fiber.New()

	// Add middleware
	app.Use(
		logger.New(), // add simple logger
	)

	// Add route
	app.Get("/send", func(c *fiber.Ctx) error {
		// Create message to publish
		message := amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(c.Query("msg")),
		}

		// Attempt to publish a message to the queue
		if err := channelRabbitMQ.Publish("", "Queue-1", false, false, message); err != nil {
			return err
		}

		return nil
	})

	// Start Fiber API server
	log.Fatal(app.Listen(":3000"))
}
