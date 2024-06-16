package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("secret_key")

type Credentials struct {
	email    string `json:"email"`
	password string `json:"password"`
}

type Claims struct {
	email string `json:"email"`
	jwt.StandardClaims
}

var client, _ = mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Register")

	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	fmt.Println(creds)
	fmt.Println(creds.email)
	fmt.Println(creds.password)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	fmt.Println("hashedPassword generated")

	collection := client.Database("weatherApp").Collection("users")
	_, err = collection.InsertOne(context.TODO(), map[string]interface{}{
		"email":    creds.email,
		"password": string(hashedPassword),
	})

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	fmt.Println("Successfully registered")

	w.WriteHeader(http.StatusCreated)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Login")

	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	fmt.Println(creds.email)
	fmt.Println(creds.password)

	var storedCreds Credentials
	collection := client.Database("weatherApp").Collection("users")
	err = collection.FindOne(context.TODO(), map[string]interface{}{"email": creds.email}).Decode(&storedCreds)
	if err != nil {
		http.Error(w, "Email not found", http.StatusUnauthorized)
		return
	}

	fmt.Println("Email found")

	if err := bcrypt.CompareHashAndPassword([]byte(storedCreds.password), []byte(creds.password)); err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	fmt.Println("Password match")

	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		email: creds.email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
}
