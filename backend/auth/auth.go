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
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

// var client, _ = mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
var client = database.Client

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")

	fmt.Println("Register")

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	fmt.Println(user)
	fmt.Println(user.Email)
	fmt.Println(user.Password)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	fmt.Println("hashedPassword generated")

	// user.Password = string(hashedPassword)
	collection := database.OpenCollection(client, "users")

	var usr models.User
	err = collection.FindOne(context.TODO(), map[string]interface{}{"email": user.Email}).Decode(&usr)
	if err == nil {
		http.Error(w, "Already exists", http.StatusBadRequest)
		return
	}

	// _, err = collection.InsertOne(context.TODO(), user)
	_, err = collection.InsertOne(context.TODO(), map[string]interface{}{
		"email":    user.Email,
		"password": string(hashedPassword),
		"city":     user.City,
		"state":    user.State,
		"country":  user.Country,
	})

	// collection := client.Database("weatherApp").Collection("users")
	// _, err = collection.InsertOne(context.TODO(), map[string]interface{}{
	// 	"email":    user.Email,
	// 	"password": string(hashedPassword),
	// })

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	fmt.Println("Successfully registered")

	res := make(map[string]bool)
	res["success"] = true
	json.NewEncoder(w).Encode(res)

	w.WriteHeader(http.StatusCreated)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Access-Control-Allow-Origin", "*")

	fmt.Println("Login")

	var creds Credentials
	// var creds map[string]string
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	fmt.Println(creds.Email)
	fmt.Println(creds.Password)

	var storedCreds Credentials
	collection := client.Database("weatherApp").Collection("users")
	err = collection.FindOne(context.TODO(), map[string]interface{}{"email": creds.Email}).Decode(&storedCreds)
	if err != nil {
		http.Error(w, "Email not found", http.StatusUnauthorized)
		return
	}

	fmt.Println("Email found")

	if err := bcrypt.CompareHashAndPassword([]byte(storedCreds.Password), []byte(creds.Password)); err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	fmt.Println("Password match")

	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Email: creds.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		fmt.Println(err)
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
