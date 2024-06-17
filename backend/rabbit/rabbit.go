package rabbit

import (
	"log"

	"github.com/streadway/amqp"
)

var conn *amqp.Connection
var ch *amqp.Channel

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func InitializeRabbitMQ(rabbitmqURL string) {
	var err error
	conn, err = amqp.Dial(rabbitmqURL)
	FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err = conn.Channel()
	FailOnError(err, "Failed to open a channel")
}
