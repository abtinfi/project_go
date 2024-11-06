# Dockerfile

# Use the official lightweight Go image based on Alpine Linux.
FROM golang:1.22-alpine

# Set the Current Working Directory inside the container to /app
WORKDIR /app

# Copy go.mod and go.sum files to the workspace
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code to the workspace
COPY . .

# Build the Go app
RUN go build -o app main.go

# Expose port 8080
EXPOSE 8080

# Command to run when starting the container
ENTRYPOINT ["./app"]
