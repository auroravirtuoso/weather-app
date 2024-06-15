module github.com/auroravirtuoso/weather-app/backend

go 1.22.2

require (
    github.com/dgrijalva/jwt-go v3.2.0+incompatible
    github.com/gorilla/mux v1.8.0
    github.com/streadway/amqp v1.0.0
    go.mongodb.org/mongo-driver v1.15.0
    golang.org/x/crypto v0.0.0-20230621151738-ffcf3497c3fe
)

replace weather-app/backend/auth => ./auth
replace weather-app/backend/weather => ./weather