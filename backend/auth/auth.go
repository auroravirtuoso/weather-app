package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/auroravirtuoso/weather-app/backend/database"
	"github.com/auroravirtuoso/weather-app/backend/models"
	"github.com/auroravirtuoso/weather-app/backend/rabbit"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// var client, _ = mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
var client = database.Client

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	collection := database.OpenCollection(client, "users")

	var usr models.User
	err = collection.FindOne(context.TODO(), map[string]interface{}{"email": user.Email}).Decode(&usr)
	if err == nil {
		http.Error(w, "Already exists", http.StatusBadRequest)
		return
	}

	user.Time = make([]string, 0)
	user.Temperature_2m = make([]string, 0)
	_, err = collection.InsertOne(context.TODO(), map[string]interface{}{
		"email":          user.Email,
		"password":       string(hashedPassword),
		"city":           user.City,
		"state":          user.State,
		"country":        user.Country,
		"time":           user.Time,
		"temperature_2m": user.Temperature_2m,
	})

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	go rabbit.ProduceWeatherData(user.Email)

	res := make(map[string]bool)
	res["success"] = true
	json.NewEncoder(w).Encode(res)

	w.WriteHeader(http.StatusCreated)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var storedCreds Credentials
	collection := client.Database("weatherApp").Collection("users")
	err = collection.FindOne(context.TODO(), map[string]interface{}{"email": creds.Email}).Decode(&storedCreds)
	if err != nil {
		http.Error(w, "Email not found", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedCreds.Password), []byte(creds.Password)); err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	var tokenExpirationTime int
	fmt.Sscanf(os.Getenv("TOKEN_EXPIRATION_TIME"), "%d", &tokenExpirationTime)
	expirationTime := time.Now().Add(time.Duration(tokenExpirationTime) * time.Minute)
	claims := &Claims{
		Email: creds.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})

	res := make(map[string]interface{})
	res["success"] = true
	res["token"] = tokenString
	json.NewEncoder(w).Encode(res)

	w.WriteHeader(http.StatusCreated)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Access-Control-Allow-Origin", "*")

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   "",
		Expires: time.Now(),
	})

	w.WriteHeader(http.StatusOK)
}

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func Authorize(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		tknStr := c.Value
		claims := &Claims{}

		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET_KEY")), nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !tkn.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		var tokenExpirationTime int
		fmt.Sscanf(os.Getenv("TOKEN_EXPIRATION_TIME"), "%d", &tokenExpirationTime)
		expirationTime := time.Now().Add(time.Duration(tokenExpirationTime) * time.Minute)
		// claims.ExpiresAt = expirationTime.Unix()
		// token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		// tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
		// if err != nil {
		// 	w.WriteHeader(http.StatusInternalServerError)
		// 	return
		// }
		http.SetCookie(w, &http.Cookie{
			Name:  "token",
			Value: tknStr,
			// Value:   tokenString,
			Expires: expirationTime,
		})

		ctx := context.WithValue(r.Context(), "email", claims.Email)
		r = r.WithContext(ctx)

		next(w, r)
	}
}
