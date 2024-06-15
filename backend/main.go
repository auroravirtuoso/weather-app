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
	r.HandleFunc("/login", auth.LoginHandler).Methods("POST")
	r.HandleFunc("/register", auth.RegisterHandler).Methods("POST")

	// Weather routes
	r.HandleFunc("/weather/{period}", weather.GetWeatherDataHandler).Methods("GET")

	log.Println("Starting server on :8080")
	http.ListenAndServe(":8080", r)
}
