package rabbit

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/auroravirtuoso/weather-app/backend/database"
	"github.com/auroravirtuoso/weather-app/backend/models"
	"github.com/auroravirtuoso/weather-app/backend/weather"
	"github.com/streadway/amqp"
)

func ProduceWeatherData(email string) {
	q, err := Ch.QueueDeclare(
		"weather_data",
		false,
		false,
		false,
		false,
		nil,
	)
	FailOnError(err, "Failed to declare a queue")

	collection := database.OpenCollection(database.Client, "users")
	var user models.User
	err = collection.FindOne(context.TODO(), map[string]interface{}{"email": email}).Decode(&user)
	if err != nil {
		FailOnError(err, "Failed to find the user")
		return
	}

	geoarr, err := weather.GetLatLonFromCity(user.City, user.State, user.Country)
	if err != nil {
		FailOnError(err, "Geocoding Error")
		return
	} else if len(geoarr) == 0 {
		FailOnError(nil, "Specified city not found")
		return
	}

	fmt.Println(geoarr)

	var start_date string
	if len(user.Time) == 0 {
		start_date = time.Now().AddDate(-3, 0, -1).Format("2001-01-01")
	} else {
		last, err := time.Parse("2001-01-01T00:00", user.Time[len(user.Time)-1])
		if err != nil {
			FailOnError(err, "Invalid Time Format")
			return
		}
		start_date = last.Format("2001-01-01")
	}

	end_date := time.Now().Format("2001-01-01")

	var api_url string = "https://archive-api.open-meteo.com/v1/era5"
	api_url += fmt.Sprintf("?latitude=%f", geoarr[0].lat)
	api_url += fmt.Sprintf("&longitude=%f", geoarr[0].lon)
	api_url += "&start_date=" + url.QueryEscape(start_date)
	api_url += "&end_date=" + url.QueryEscape(end_date)
	api_url += "&hourly=" + "temperature_2m"
	// api_url = url.QueryEscape(api_url)
	fmt.Println(api_url)
	resp, err := http.Get(api_url)
	if err != nil {
		FailOnError(err, "API Error")
	}
	defer resp.Body.Close()

	fmt.Println(resp.Body)

	var body map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		FailOnError(err, "JSON Error")
	}

	body_hourly := body["hourly"].(map[string]interface{})

	time_arr := body_hourly["time"].([]interface{})
	temperature_2m := body_hourly["temperature_2m"].([]interface{})

	var cnt int = 0
	for i := 0; i < len(temperature_2m); i++ {
		val, err := json.Marshal(temperature_2m[i])
		if err != nil {
			FailOnError(err, "Marshal Failure")
			break
		} else if string(val) == "\"null\"" {
			break
		}
		cnt++
	}
	time_arr = time_arr[:cnt]
	temperature_2m = temperature_2m[:cnt]

	messageBody := map[string]interface{}{
		"email": user.Email,
		"data": map[string]interface{}{
			"time":           time_arr,
			"temperature_2m": temperature_2m,
		},
	}
	var messageJSON []byte
	messageJSON, err = json.Marshal(messageBody)

	err = Ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(messageJSON),
		})
	FailOnError(err, "Failed to publish a message")

	log.Printf(" [x] Sent %s", body)
}
