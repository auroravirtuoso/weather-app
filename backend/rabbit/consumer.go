package rabbit

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/auroravirtuoso/weather-app/backend/database"
	"github.com/auroravirtuoso/weather-app/backend/models"
	"go.mongodb.org/mongo-driver/bson"
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
			log.Println("Received a message")
			// log.Printf("Received a message: %s", d.Body)
			// Process and store data in MongoDB
			var body map[string]interface{}
			err := json.Unmarshal(d.Body, &body)
			if err != nil {
				FailOnError(err, "Unmarshal failure")
				continue
			}
			email := body["email"].(string)
			data := body["data"].(map[string]interface{})
			time_arr := data["time"].([]interface{})
			temperature_2m := data["temperature_2m"].([]interface{})

			collection := database.OpenCollection(database.Client, "users")
			var user models.User
			err = collection.FindOne(context.TODO(), map[string]interface{}{"email": email}).Decode(&user)
			if err != nil {
				FailOnError(err, "User not found")
				continue
			}

			var idx int = 0
			if len(user.Time) > 0 {
				last, err := time.Parse("2006-01-02T15:04", user.Time[len(user.Time)-1])
				if err != nil {
					FailOnError(err, "Invalid Time Format")
					break
				}
				for ; idx < len(time_arr); idx++ {
					cur, err := time.Parse("2006-01-02T15:04", time_arr[idx].(string))
					if err != nil {
						FailOnError(err, "Invalid Time Format")
						break
					}
					if last.Before(cur) {
						break
					}
				}
			}

			for ; idx < len(time_arr); idx++ {
				user.Time = append(user.Time, time_arr[idx].(string))
				user.Temperature_2m = append(user.Temperature_2m, fmt.Sprintf("%f", temperature_2m[idx].(float64)))
			}

			filter := bson.M{"email": user.Email}
			update := bson.M{
				"$set": bson.M{
					"time":           user.Time,
					"temperature_2m": user.Temperature_2m,
				},
			}

			result, err := collection.UpdateOne(context.TODO(), filter, update)
			if err != nil {
				log.Fatal(err)
			}

			log.Printf("Matched %v documents and updated %v documents.\n", result.MatchedCount, result.ModifiedCount)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
