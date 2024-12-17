# golang1.20 as the base image
FROM golang:1.23.2 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

# Build the Go app
RUN go build -o main cmd/server/main.go

# Start a new stage from the alpine image
FROM alpine:latest  

# Install necessary libraries for MySQL client (e.g., for database connection)
RUN apk --no-cache add ca-certificates

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the binary from the previous stage
COPY --from=builder /app/main .

# Expose the port the app runs on
EXPOSE 50051

# Command to run the executable
CMD ["./main"]