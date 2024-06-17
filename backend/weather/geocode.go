package weather

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

type Geolocation struct {
	lat float64
	lon float64
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
	api_url += "&appid=9bd398148984a3f361fa58d491cc53e5" // + os.Getenv("OPENWEATHERMAP_API_KEY")
	// api_url = url.QueryEscape(api_url)
	fmt.Println(api_url)
	resp, e := http.Get(api_url)
	if e != nil {
		err = e
		return
	}
	defer resp.Body.Close()

	fmt.Println("BODY")
	fmt.Println(resp.Body)

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
		geo.lat = result["lat"].(float64)
		geo.lon = result["lon"].(float64)
		geoarr = append(geoarr, geo)
	}

	return
}

func GetLatLonFromCityHandler(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Access-Control-Allow-Origin", "*")

	fmt.Println("GetLatLonFromCityHandler")
	var vars map[string]string
	err := json.NewDecoder(r.Body).Decode(&vars)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	city := vars["city"]
	state := vars["state"]
	country := vars["country"]
	fmt.Println(city)
	fmt.Println(state)
	fmt.Println(country)
	geoarr, err := GetLatLonFromCity(city, state, country)
	if err != nil {
		http.Error(w, "Fetch Error", http.StatusInternalServerError)
	}

	fmt.Println(geoarr)

	if len(geoarr) == 0 {
		http.Error(w, "Specified city not found", http.StatusNotFound)
	}

	results := make(map[string]float64)
	results["lat"] = geoarr[0].lat
	results["lon"] = geoarr[0].lon
	json.NewEncoder(w).Encode(results)
}
