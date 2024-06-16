package weather

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/auroravirtuoso/weather-app/backend/auth"
	"github.com/dgrijalva/jwt-go"
)

// var client, _ = mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))

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

	tokenCookie, err := r.Cookie("token")
	if err != nil {
		http.Error(w, "Token not found", http.StatusUnauthorized)
	}

	/* Authentication */
	// tokenStr := r.Header.Get("Authorization")
	tokenStr := tokenCookie.Value
	claims := &auth.Claims{}

	tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret_key"), nil
	})

	if err != nil || !tkn.Valid {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	geoarr, err := GetLatLonFromCity(city, state, country)
	if err != nil {
		http.Error(w, "Geocoding Error", http.StatusInternalServerError)
	} else if len(geoarr) == 0 {
		http.Error(w, "Specified city not found", http.StatusNotFound)
	}

	fmt.Println(geoarr)

	var api_url string = "https://archive-api.open-meteo.com/v1/era5"
	api_url += fmt.Sprintf("?latitude=%f", geoarr[0].lat)
	api_url += fmt.Sprintf("&longitude=%f", geoarr[0].lon)
	api_url += "&start_date=" + url.QueryEscape(start_date)
	api_url += "&end_date=" + url.QueryEscape(end_date)
	api_url += "&hourly=" + url.QueryEscape(hourly)
	// api_url = url.QueryEscape(api_url)
	fmt.Println(api_url)
	resp, err := http.Get(api_url)
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
