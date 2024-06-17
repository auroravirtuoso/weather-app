package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/auroravirtuoso/weather-app/backend/auth"
	"github.com/joho/godotenv"
	"github.com/streadway/amqp"

	"github.com/auroravirtuoso/weather-app/backend/weather"

	"github.com/gorilla/mux"
)

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow requests from any origin
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Allow specified HTTP methods
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

		// Allow specified headers
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")

		// Continue with the next handler
		next.ServeHTTP(w, r)
	})
}

func CORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", os.Getenv("REACT_APP_FRONTEND_URL"))
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

		if r.Method == "OPTIONS" {
			http.Error(w, "No Content", http.StatusNoContent)
			return
		}

		next(w, r)
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func produceWeatherData(ch *amqp.Channel) {
	q, err := ch.QueueDeclare(
		"weather_data",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

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
	failOnError(err, "Failed to publish a message")

	log.Printf(" [x] Sent %s", body)
}

func consumeWeatherData(ch *amqp.Channel) {
	q, err := ch.QueueDeclare(
		"weather_data",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register a consumer")

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

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	rabbitmqURL := os.Getenv("RABBITMQ_URL")

	fmt.Println(rabbitmqURL)

	conn, err := amqp.Dial(rabbitmqURL)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// ch, err := conn.Channel()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// produceWeatherData(ch)
	// consumeWeatherData(ch)

	r := mux.NewRouter()

	// Auth routes

	// r.Use(enableCORS)
	r.HandleFunc("/api/v1/login", CORS(auth.LoginHandler)).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/v1/register", CORS(auth.RegisterHandler)).Methods("POST", "OPTIONS")

	// Weather routes
	r.HandleFunc("/api/v1/weather", CORS(weather.GetWeatherDataHandler)).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/v1/geocode", CORS(weather.GetLatLonFromCityHandler)).Methods("POST", "OPTIONS")

	log.Println("Starting server on :8080")
	http.ListenAndServe(":8080", r)
}
