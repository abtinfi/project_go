# Stage 1: Build
FROM golang:1.22-alpine AS builder

# Install required build tools
RUN apk add --no-cache git

# Set the Current Working Directory inside the container to /app
WORKDIR /app

# Copy go.mod and go.sum files to the workspace
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire source code to the workspace
COPY . .

# Build the Go app
RUN go build -o app main.go

# Stage 2: Run
FROM alpine:latest

# Set the working directory
WORKDIR /root/

# Copy the built binary from the builder
COPY --from=builder /app/app .

# Expose port 8080
EXPOSE 8080

# Command to run the executable
ENTRYPOINT ["./app"]
