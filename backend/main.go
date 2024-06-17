package main

import (
	"log"
	"net/http"
	"os"

	"github.com/auroravirtuoso/weather-app/backend/auth"
	"github.com/auroravirtuoso/weather-app/backend/middlewares"
	"github.com/auroravirtuoso/weather-app/backend/rabbit"
	"github.com/joho/godotenv"

	"github.com/auroravirtuoso/weather-app/backend/weather"

	"github.com/gorilla/mux"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	rabbitmqURL := os.Getenv("RABBITMQ_URL")
	rabbit.InitializeRabbitMQ(rabbitmqURL)

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
	r.HandleFunc("/api/v1/login", middlewares.CORS(auth.LoginHandler)).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/v1/register", middlewares.CORS(auth.RegisterHandler)).Methods("POST", "OPTIONS")

	// Weather routes
	r.HandleFunc("/api/v1/weather", middlewares.CORS(weather.GetWeatherDataHandler)).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/v1/geocode", middlewares.CORS(weather.GetLatLonFromCityHandler)).Methods("POST", "OPTIONS")

	log.Println("Starting server on :8080")
	http.ListenAndServe(":8080", r)
}
