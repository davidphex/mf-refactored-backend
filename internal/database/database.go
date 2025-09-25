package database

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var Client *mongo.Client

func ConnectToDB() (*mongo.Client, error) {
	MONGO_URI := os.Getenv("MONGO_URI")

	if MONGO_URI == "" {
		log.Fatal("MONGO_URI environment variable not set")
	}

	log.Println("Connecting to MongoDB...")

	// Uses the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	// Defines the options for the MongoDB client
	opts := options.Client().ApplyURI(MONGO_URI).SetServerAPIOptions(serverAPI)
	// Creates a new client and connects to the server
	client, err := mongo.Connect(opts)

	if err != nil {
		panic(err)
	}
	/*defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()*/

	log.Println("Connected to MongoDB")

	return client, nil
}

func DisconnectDB(client *mongo.Client) {
	if err := client.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
}
