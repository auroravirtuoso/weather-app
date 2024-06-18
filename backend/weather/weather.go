package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/auroravirtuoso/weather-app/backend/database"
	"github.com/auroravirtuoso/weather-app/backend/geolocation"
	"github.com/auroravirtuoso/weather-app/backend/models"
)

// var client, _ = mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))

// https://open-meteo.com/en/docs
func GetWeatherDataHandler(w http.ResponseWriter, r *http.Request) {
	email := r.Context().Value("email").(string)

	var user models.User
	collection := database.OpenCollection(database.Client, "users")
	err := collection.FindOne(context.TODO(), map[string]interface{}{"email": email}).Decode(&user)
	if err != nil {
		http.Error(w, "Specified email not found", http.StatusInternalServerError)
		return
	}

	var city string
	query := r.URL.Query()
	if query.Has("city") {
		city = query.Get("city")
	} else {
		city = user.City
		// http.Error(w, "city is required", http.StatusBadRequest)
	}
	var state string
	if query.Has("state") {
		state = query.Get("state")
	} else {
		state = user.State
	}
	var country string
	if query.Has("country") {
		country = query.Get("country")
	} else {
		country = user.Country
	}
	var start_date string
	if query.Has("start_date") {
		start_date = query.Get("start_date")
	} else {
		http.Error(w, "start_date is required", http.StatusBadRequest)
	}
	var end_date string
	if query.Has("end_date") {
		end_date = query.Get("end_date")
	} else {
		http.Error(w, "end_date is required", http.StatusBadRequest)
	}
	var hourly string
	if query.Has("hourly") {
		hourly = query.Get("hourly")
	} else {
		hourly = "temperature_2m"
		// http.Error(w, "hourly is required", http.StatusBadRequest)
	}

	geoarr, err := geolocation.GetLatLonFromCity(city, state, country)
	if err != nil {
		http.Error(w, "Geocoding Error", http.StatusInternalServerError)
	} else if len(geoarr) == 0 {
		http.Error(w, "Specified city not found", http.StatusNotFound)
	}

	var api_url string = "https://archive-api.open-meteo.com/v1/era5"
	api_url += fmt.Sprintf("?latitude=%f", geoarr[0].Lat)
	api_url += fmt.Sprintf("&longitude=%f", geoarr[0].Lon)
	api_url += "&start_date=" + url.QueryEscape(start_date)
	api_url += "&end_date=" + url.QueryEscape(end_date)
	api_url += "&hourly=" + url.QueryEscape(hourly)
	resp, err := http.Get(api_url)
	if err != nil {
		http.Error(w, "API Error", http.StatusInternalServerError)
	}
	defer resp.Body.Close()

	var body map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		http.Error(w, "JSON Error", http.StatusInternalServerError)
	}

	body_hourly := body["hourly"].(map[string]interface{})

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"results": body_hourly,
	})
}
