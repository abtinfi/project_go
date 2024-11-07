// main.go
package main

import (
	"log"
	"os"
	"project_go/api"
	"project_go/config"

	"github.com/joho/godotenv"
	"github.com/kataras/iris/v12"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Could not load .env file, using environment variables instead")
	}

	// Initialize Cassandra connection using environment variables
	config.InitCassandra()
	defer config.CloseSession()

	// Create a new Iris application
	app := iris.New()

	// Define routes
	app.Post("/users", api.CreateUser)
	app.Get("/users/{id}", api.GetUser)
	app.Put("/users/{id}", api.UpdateUser)
	app.Delete("/users/{id}", api.DeleteUser)
	app.Get("/users", api.ListUsers)

	// Start the application on the port defined in the environment variable
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080" // Default port if not set
		log.Printf("Using default port %s as APP_PORT is not set", port)
	}

	log.Printf("Starting server on port %s...", port)
	err = app.Listen(":" + port)
	if err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}
