# syntax=docker/dockerfile:1
FROM golang:1.23.1 AS build

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o /todo cmd/todo_app/main.go

# Start a new stage from scratch
FROM alpine:latest

WORKDIR /

# Copy the Pre-built binary file from the previous stage
COPY --from=build /todo /todo

# Command to run the executable
CMD ["/todo"]
