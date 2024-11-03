# Dockerfile
FROM golang:1.22-alpine

# Set up the working directory
WORKDIR /app

# Enable Go modules (no need to set GOPATH for module-based builds)
ENV GO111MODULE=on

# Copy go.mod and go.sum first to cache dependency resolution
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application files
COPY . .

# Build the application
RUN go build -o app main.go

# Command to run the application
CMD ["./app"]
