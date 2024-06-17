package rabbit

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/auroravirtuoso/weather-app/backend/database"
	"github.com/auroravirtuoso/weather-app/backend/models"
)

func ConsumeWeatherData() {
	q, err := Ch.QueueDeclare(
		"weather_data",
		false,
		false,
		false,
		false,
		nil,
	)
	FailOnError(err, "Failed to declare a queue")

	msgs, err := Ch.Consume(
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
			var body map[string]interface{}
			err := json.Unmarshal(d.Body, &body)
			if err != nil {
				FailOnError(err, "Unmarshal failure")
				continue
			}
			email := body["email"].(string)
			data := body["data"].(map[string]interface{})
			time_arr := data["time"].([]string)
			temperature_2m := data["temperature_2m"].([]float64)

			collection := database.OpenCollection(database.Client, "users")
			var user models.User
			err = collection.FindOne(context.TODO(), map[string]interface{}{"email": email}).Decode(&user)
			if err != nil {
				FailOnError(err, "User not found")
				continue
			}

			var idx int = 0
			if len(user.Time) > 0 {
				last, err := time.Parse("2001-01-01T00:00", user.Time[len(user.Time)-1])
				if err != nil {
					FailOnError(err, "Invalid Time Format")
					break
				}
				for ; idx < len(time_arr); idx++ {
					cur, err := time.Parse("2001-01-01T00:00", user.Time[len(user.Time)-1])
					if err != nil {
						FailOnError(err, "Invalid Time Format")
						break
					}
					if last.Before(cur) {
						break
					}
				}
			}
			user.Time = append(user.Time, time_arr[idx:]...)
			user.Temperature = append(user.Temperature, temperature_2m[idx:]...)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
