package main

import (
	"log"
	"os"
	"strconv"

	"github.com/davidphex/memoryframe-backend/internal/app"
	"github.com/davidphex/memoryframe-backend/internal/database"
)

func main() {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	portInt, err := strconv.Atoi(port)
	if err != nil {
		log.Fatalf("Invalid port number: %v", err)
	}

	cfg := app.Config{
		Port: portInt,
		Env:  "development",
	}

	client, err := database.ConnectToDB()
	if err != nil {
		log.Fatal(err)
	}
	// Disconnect the client upon exit
	defer database.DisconnectDB(client)

	application := app.New(cfg, client)

	if err := application.Serve(); err != nil {
		log.Fatal(err)
	}
}
