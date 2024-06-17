package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/auroravirtuoso/weather-app/backend/auth"
	"github.com/auroravirtuoso/weather-app/backend/database"
	"github.com/auroravirtuoso/weather-app/backend/models"
	"github.com/auroravirtuoso/weather-app/backend/rabbit"
	"github.com/dgrijalva/jwt-go"
)

// https://open-meteo.com/en/docs
func GetUserWeatherDataHandler(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Access-Control-Allow-Origin", "*")
	// vars := mux.Vars(r)

	tokenCookie, err := r.Cookie("token")
	if err != nil {
		http.Error(w, "Token not found", http.StatusUnauthorized)
	}

	/* Authentication */
	// tokenStr := r.Header.Get("Authorization")
	tokenStr := tokenCookie.Value
	fmt.Println(tokenStr)
	claims := &auth.Claims{}

	tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	if err != nil || !tkn.Valid {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	fmt.Println("Authorized")

	var user models.User
	collection := database.OpenCollection(database.Client, "users")
	err = collection.FindOne(context.TODO(), map[string]interface{}{"email": claims.Email}).Decode(&user)
	if err != nil {
		http.Error(w, "Specified email not found", http.StatusInternalServerError)
		return
	}

	go rabbit.ProduceWeatherData()

	results := make(map[string]interface{})
	results["time"] = user.Time
	results["temperature_2m"] = user.Temperature

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"results": results,
	})
}
