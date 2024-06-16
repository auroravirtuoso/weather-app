package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client, _ = mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))

type WeatherData struct {
	Date time.Time `json:"date"`
	Temp float64   `json:"temp"`
}

type Geolocation struct {
	lat float64
	lon float64
}

// https://openweathermap.org/api/geocoding-api
func GetLatLonFromCity(city string, state string, country string) (geoarr []Geolocation, err error) {
	geoarr = make([]Geolocation, 0)

	var url string = "https://api.openweathermap.org/geo/1.0/direct?q="
	url += city
	if len(state) > 0 {
		url += ","
		url += state
	}
	if len(country) > 0 {
		url += ","
		url += country
	}
	url += "&limit=5"
	url += "&appid=9bd398148984a3f361fa58d491cc53e5" // + os.Getenv("OPENWEATHERMAP_API_KEY")
	fmt.Println(url)
	resp, e := http.Get(url)
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

	if len(geoarr) == 0 {
		http.Error(w, "Specified city not found", http.StatusNotFound)
	}

	results := make(map[string]float64)
	results["lat"] = geoarr[0].lat
	results["lon"] = geoarr[0].lon
	json.NewEncoder(w).Encode(results)
}

type HourlyUnits struct {
	time          string `json:"time"`
	temperature2m string `json:"temperature_2m"`
}

type WeatherResponse struct {
	time          []string  `json:"time"`
	temperature2m []float64 `json:"temperature_2m"`
}

type OpenMeteoHistoryBody struct {
	latitude              float64         `json:"latitude"`
	longitude             float64         `json:"longitude"`
	generationtime_ms     float64         `json:"generationtime_ms"`
	utc_offset_seconds    int             `json:"utc_offset_seconds"`
	timezone              string          `json:"timezone"`
	timezone_abbreviation string          `json:"timezone_abbreviation"`
	elevation             float64         `json:"elevation"`
	hourly_units          HourlyUnits     `json:"hourly_units"`
	hourly                WeatherResponse `json:"hourly"`
}

// https://open-meteo.com/en/docs
func GetWeatherDataHandler(w http.ResponseWriter, r *http.Request) {
	// vars := mux.Vars(r)

	var city string
	query := r.URL.Query()
	fmt.Println(query)
	if query.Has("city") {
		city = query.Get("city")
	} else {
		http.Error(w, "city is required", http.StatusBadRequest)
	}
	var state string
	if query.Has("state") {
		state = query.Get("state")
	}
	var country string
	if query.Has("country") {
		country = query.Get("country")
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
		http.Error(w, "hourly is required", http.StatusBadRequest)
	}

	// fmt.Println("----------")
	// fmt.Println(city)
	// fmt.Println(state)
	// fmt.Println(country)
	// fmt.Println(start_date)
	// fmt.Println(end_date)
	// fmt.Println(hourly)
	// fmt.Println("----------")

	/* Authentication */
	// tokenStr := r.Header.Get("Authorization")
	// claims := &auth.Claims{}

	// tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
	// 	return []byte("secret_key"), nil
	// })

	// if err != nil || !tkn.Valid {
	// 	http.Error(w, "Unauthorized", http.StatusUnauthorized)
	// 	return
	// }

	geoarr, err := GetLatLonFromCity(city, state, country)
	if err != nil {
		http.Error(w, "Geocoding Error", http.StatusInternalServerError)
	} else if len(geoarr) == 0 {
		http.Error(w, "Specified city not found", http.StatusNotFound)
	}

	fmt.Println(geoarr)

	var url string = "https://archive-api.open-meteo.com/v1/era5"
	url += fmt.Sprintf("?latitude=%f", geoarr[0].lat)
	url += fmt.Sprintf("&longitude=%f", geoarr[0].lon)
	url += "&start_date=" + start_date
	url += "&end_date=" + end_date
	url += "&hourly=" + hourly
	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "API Error", http.StatusInternalServerError)
	}
	defer resp.Body.Close()

	fmt.Println(resp.Body)

	var body map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		http.Error(w, "JSON Error", http.StatusInternalServerError)
	}

	fmt.Println(body["latitude"].(float64))
	fmt.Println(body["longitude"].(float64))

	body_hourly := body["hourly"].(map[string]interface{})

	// time := body_hourly["time"].([]string)
	// temperature_2m := body_hourly["temperature_2m"].([]float64)

	json.NewEncoder(w).Encode(body_hourly)

	// collection := client.Database("weatherApp").Collection("weatherData")

	// var results []WeatherData
	// var cursor *mongo.Cursor

	// switch period {
	// case "month":
	// 	cursor, err = collection.Find(context.TODO(), map[string]interface{}{"date": map[string]interface{}{"$gte": time.Now().AddDate(0, -1, 0)}})
	// case "year":
	// 	cursor, err = collection.Find(context.TODO(), map[string]interface{}{"date": map[string]interface{}{"$gte": time.Now().AddDate(-1, 0, 0)}})
	// case "3years":
	// 	cursor, err = collection.Find(context.TODO(), map[string]interface{}{"date": map[string]interface{}{"$gte": time.Now().AddDate(-3, 0, 0)}})
	// default:
	// 	http.Error(w, "Invalid period", http.StatusBadRequest)
	// 	return
	// }

	// if err != nil {
	// 	http.Error(w, "Internal server error", http.StatusInternalServerError)
	// 	return
	// }

	// defer cursor.Close(context.TODO())
	// for cursor.Next(context.TODO()) {
	// 	var weather WeatherData
	// 	err := cursor.Decode(&weather)
	// 	if err != nil {
	// 		http.Error(w, "Internal server error", http.StatusInternalServerError)
	// 		return
	// 	}
	// 	results = append(results, weather)
	// }

	// json.NewEncoder(w).Encode(results)
}
