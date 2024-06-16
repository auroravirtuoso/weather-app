package main

import (
	"log"
	"net/http"

	"github.com/auroravirtuoso/weather-app/backend/auth"

	"github.com/auroravirtuoso/weather-app/backend/weather"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	// Auth routes
	r.HandleFunc("/api/v1/login", auth.LoginHandler).Methods("POST")
	r.HandleFunc("/api/v1/register", auth.RegisterHandler).Methods("POST")

	// Weather routes
	r.HandleFunc("/api/v1/weather", weather.GetWeatherDataHandler).Methods("GET")
	r.HandleFunc("/api/v1/geocode", weather.GetLatLonFromCityHandler).Methods("POST")

	log.Println("Starting server on :8080")
	http.ListenAndServe(":8080", r)
}
