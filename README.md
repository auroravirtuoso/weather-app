# Weather-App
Welcome to the Weather-App project! This repository contains a weather application built using GoLang, MongoDB, RabbitMQ for the backend, and React for the frontend.
## Overview
Weather-App is a simple application that allows users to register, log in, and view historical weather data for their hometown. The backend is built with GoLang, MongoDB, and RabbitMQ, while the frontend is built with React.
## Features
- User registration and authentication
- Fetch historical weather data using OpenWeatherMap and Open-Meteo APIs
- Real-time data processing with RabbitMQ
- Secure API endpoints with JWT-based authentication
## Prerequisites
Before you begin, ensure you have the following installed on your local machine:
- Go 1.16+
- MongoDB
- RabbitMQ
- Docker (optional, for containerized deployment)
## Installation
### 1. Clone the repository:
```console
git clone https://github.com/auroravirtuoso/weather-app.git
cd weather-app
```
### 2. Set up the backend:
Navigate to the backend directory and install the dependencies:
```console
cd backend
go mod tidy
```
### 3. Set up the frontend:
Navigate to the frontend directory and install the dependencies:
```console
cd frontend
npm install
```
## Environment Variables
Set the environment variables based on .env.example file.
### backend/.env
```
MONGODB_URI=mongodb://mongo:27017
MONGODB_DATABASE=weatherApp
OPENWEATHERMAP_API_KEY=YOUR_API_KEY
RABBITMQ_URL=amqp://guest:guest@rabbitmq:5672/
REACT_APP_FRONTEND_URL=http://localhost
JWT_SECRET_KEY=secret_key
TOKEN_EXPIRATION_TIME=5
```
### frontend/.env
```
REACT_APP_BACKEND_URL=http://localhost:8080
```
## Usage
Run the following command in terminal.
```sh
sudo docker-compose up --build
```
## API Endpoints
### Authentication
- **POST /register**: Register a new user
- **POST /login**: Log in a user and receive a JWT token
- **POST /logout**: Log out a user
### Weather Data
- **GET /userweather**: Get historical weather data for the authenticated user
- **GET /weather/city**: Get weather data for a specific city (requires query parameters: `city`, `state`, `country`, `start_date`, `end_date`, `hourly`)
### Geolocation
- **GET /geolocation**: Get latitude and longitude for a specific city (requires query parameters: `city`, `state`, `country`)
## License
This project is licensed under the MIT License. See the LICENSE file for more information.