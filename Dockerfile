# Use the official Go image as a base image
FROM golang:1.21.6 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules files
COPY go.mod .
COPY go.sum .

# Download dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Disable CGO
ENV CGO_ENABLED=0

# Build the executable
RUN go build -o db-syncer main.go

# Start a new stage from scratch
FROM alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the executable from the builder stage to the new image
COPY --from=builder /app/db-syncer /app

# Expose the port your application listens on
EXPOSE 80

# Command to run the executable
CMD ["./db-syncer"]
