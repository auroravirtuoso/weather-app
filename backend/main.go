package main

import (
	"log"
	"net/http"
	"os"

	"github.com/auroravirtuoso/weather-app/backend/auth"
	"github.com/auroravirtuoso/weather-app/backend/geolocation"
	"github.com/auroravirtuoso/weather-app/backend/middlewares"
	"github.com/auroravirtuoso/weather-app/backend/rabbit"
	"github.com/joho/godotenv"

	"github.com/auroravirtuoso/weather-app/backend/weather"

	"github.com/gorilla/mux"
)

func main() {
	// load .env
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// initialize RabbitMQ
	rabbit.InitializeRabbitMQ(os.Getenv("RABBITMQ_URL"))
	defer rabbit.Conn.Close()

	go rabbit.ConsumeWeatherData()

	r := mux.NewRouter()

	// Auth routes
	r.HandleFunc("/api/v1/login", middlewares.CORS(auth.LoginHandler)).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/v1/logout", middlewares.CORS(auth.LogoutHandler)).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/v1/register", middlewares.CORS(auth.RegisterHandler)).Methods("POST", "OPTIONS")

	// Weather routes
	r.HandleFunc("/api/v1/weather", middlewares.CORS(auth.Authorize(weather.GetWeatherDataHandler))).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/v1/userweather", middlewares.CORS(auth.Authorize(weather.GetUserWeatherDataHandler))).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/v1/geocode", middlewares.CORS(auth.Authorize(geolocation.GetLatLonFromCityHandler))).Methods("GET", "OPTIONS")

	log.Println("Starting server on :8080")
	http.ListenAndServe(":8080", r)
}
