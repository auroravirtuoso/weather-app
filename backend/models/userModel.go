package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// User is the model that governs all notes objects retrived or inserted into the DB
type User struct {
	ID             primitive.ObjectID `bson:"_id"`
	Email          string             `json:"email" validate:"email,required"`
	Password       string             `json:"password" validate:"required,min=6""`
	City           string             `json:"city" validate:"required"`
	State          string             `json:"state"`
	Country        string             `json:"country" validate:"required"`
	Time           []string           `json:"time"`
	Temperature_2m []string           `json:"temperature_2m"`
	// Token         string             `json:"token"`
	// Refresh_token string             `json:"refresh_token"`
	// Created_at time.Time `json:"created_at"`
	// Updated_at time.Time `json:"updated_at"`
	// User_id       string             `json:"user_id"`
}
