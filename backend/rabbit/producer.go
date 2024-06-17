package rabbit

import (
	"log"

	"github.com/streadway/amqp"
)

func produceWeatherData() {
	q, err := ch.QueueDeclare(
		"weather_data",
		false,
		false,
		false,
		false,
		nil,
	)
	FailOnError(err, "Failed to declare a queue")

	body := "Weather data payload"
	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	FailOnError(err, "Failed to publish a message")

	log.Printf(" [x] Sent %s", body)
}
