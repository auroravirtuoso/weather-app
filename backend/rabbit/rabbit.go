package rabbit

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

var Conn *amqp.Connection
var Ch *amqp.Channel

func InitializeRabbitMQ(rabbitmqURL string) {
	fmt.Println(rabbitmqURL)
	var err error
	Conn, err = amqp.Dial(rabbitmqURL)
	FailOnError(err, "Failed to connect to RabbitMQ")

	Ch, err = Conn.Channel()
	FailOnError(err, "Failed to open a channel")
}

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
