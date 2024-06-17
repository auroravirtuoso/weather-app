package rabbit

import (
	"log"
)

func consumeWeatherData() {
	q, err := ch.QueueDeclare(
		"weather_data",
		false,
		false,
		false,
		false,
		nil,
	)
	FailOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	FailOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			// Process and store data in MongoDB
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
