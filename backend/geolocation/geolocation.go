package geolocation

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

type Geolocation struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

// https://openweathermap.org/api/geocoding-api
func GetLatLonFromCity(city string, state string, country string) (geoarr []Geolocation, err error) {
	geoarr = make([]Geolocation, 0)

	var api_url string = "https://api.openweathermap.org/geo/1.0/direct?q="
	api_url += url.QueryEscape(city)
	if len(state) > 0 {
		api_url += ","
		api_url += url.QueryEscape(state)
	}
	if len(country) > 0 {
		api_url += ","
		api_url += url.QueryEscape(country)
	}
	api_url += "&limit=5"
	api_url += "&appid=" + os.Getenv("OPENWEATHERMAP_API_KEY")
	client := http.Client{
		Timeout: 1 * time.Minute,
	}
	resp, e := client.Get(api_url)
	if e != nil {
		log.Println("GEOAPI ERROR")
		err = e
		return
	}
	defer resp.Body.Close()

	var data []interface{}
	e = json.NewDecoder(resp.Body).Decode(&data)
	if e != nil {
		err = e
		log.Fatalf("Failed to decode geolocation data: %v", e)
		return
	}

	for _, itf := range data {
		result := itf.(map[string]interface{})
		var geo Geolocation
		geo.Lat = result["lat"].(float64)
		geo.Lon = result["lon"].(float64)
		geoarr = append(geoarr, geo)
	}

	return
}

func GetLatLonFromCityHandler(w http.ResponseWriter, r *http.Request) {
	var vars map[string]string
	err := json.NewDecoder(r.Body).Decode(&vars)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	query := r.URL.Query()
	var city string
	if query.Has("city") {
		city = query.Get("city")
	} else {
		http.Error(w, "city is required", http.StatusBadRequest)
	}
	var state string
	if query.Has("state") {
		state = query.Get("state")
	} else {
		http.Error(w, "state is required", http.StatusBadRequest)
	}
	var country string
	if query.Has("country") {
		country = query.Get("country")
	} else {
		http.Error(w, "country is required", http.StatusBadRequest)
	}

	geoarr, err := GetLatLonFromCity(city, state, country)
	if err != nil {
		http.Error(w, "Fetch Error", http.StatusInternalServerError)
	}

	if len(geoarr) == 0 {
		http.Error(w, "Specified city not found", http.StatusNotFound)
	}

	results := make(map[string]float64)
	results["lat"] = geoarr[0].Lat
	results["lon"] = geoarr[0].Lon
	json.NewEncoder(w).Encode(results)
}
