FROM golang:1.22

# Set work directory
WORKDIR /usr/src/app

# Copy all source
COPY . .

# Download dependencies
RUN go mod download 

# Build module
RUN go build -o main .

# Set port
EXPOSE 8080

# Launch backend
CMD ["./main"]