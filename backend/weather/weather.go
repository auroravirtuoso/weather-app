package weather

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/auroravirtuoso/weather-app/backend/auth"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client, _ = mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://mongo:27017"))

type WeatherData struct {
	Date time.Time `json:"date"`
	Temp float64   `json:"temp"`
}

func GetWeatherDataHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	period := vars["period"]

	tokenStr := r.Header.Get("Authorization")
	claims := &auth.Claims{}

	tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret_key"), nil
	})

	if err != nil || !tkn.Valid {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	collection := client.Database("weatherApp").Collection("weatherData")

	var results []WeatherData
	var cursor *mongo.Cursor

	switch period {
	case "month":
		cursor, err = collection.Find(context.TODO(), map[string]interface{}{"date": map[string]interface{}{"$gte": time.Now().AddDate(0, -1, 0)}})
	case "year":
		cursor, err = collection.Find(context.TODO(), map[string]interface{}{"date": map[string]interface{}{"$gte": time.Now().AddDate(-1, 0, 0)}})
	case "3years":
		cursor, err = collection.Find(context.TODO(), map[string]interface{}{"date": map[string]interface{}{"$gte": time.Now().AddDate(-3, 0, 0)}})
	default:
		http.Error(w, "Invalid period", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		var weather WeatherData
		err := cursor.Decode(&weather)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		results = append(results, weather)
	}

	json.NewEncoder(w).Encode(results)
}
