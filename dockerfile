# Start from the official Golang image with Alpine
FROM golang:1.20-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o library-api-with-jwt .

# Start a new stage from scratch
FROM alpine:latest

# Set the working directory
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/library-api-with-jwt .

# Copy the .env file into the container
COPY .env ./

# Expose the port (default is 8080)
EXPOSE 8080

# Command to run the application
CMD ["./library-api-with-jwt"]
