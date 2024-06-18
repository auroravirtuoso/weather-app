package weather

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/auroravirtuoso/weather-app/backend/database"
	"github.com/auroravirtuoso/weather-app/backend/models"
	"github.com/auroravirtuoso/weather-app/backend/rabbit"
)

// https://open-meteo.com/en/docs
func GetUserWeatherDataHandler(w http.ResponseWriter, r *http.Request) {
	email := r.Context().Value("email").(string)

	var user models.User
	collection := database.OpenCollection(database.Client, "users")
	err := collection.FindOne(context.TODO(), map[string]interface{}{"email": email}).Decode(&user)
	if err != nil {
		http.Error(w, "Specified email not found", http.StatusInternalServerError)
		return
	}

	go rabbit.ProduceWeatherData(email)

	results := make(map[string]interface{})
	results["time"] = user.Time
	results["temperature_2m"] = user.Temperature_2m

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"results": results,
	})
}
