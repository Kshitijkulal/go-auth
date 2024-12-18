package main

import (
    "context"
    "log"
    "os"

    "github.com/joho/godotenv"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func init() {
    // Load .env file
    err := godotenv.Load()
    if err != nil {
        log.Println("No .env file found, using system environment variables")
    }
}

func ConnectDB() {
    // Retrieve MongoDB URI from environment variables
    mongoURI := os.Getenv("MONGO_URI")
    if mongoURI == "" {
        log.Fatal("MONGO_URI environment variable is not set")
    }

    // Set client options
    clientOptions := options.Client().ApplyURI(mongoURI)

    // Connect to MongoDB
    var err error
    client, err = mongo.Connect(context.TODO(), clientOptions)
    if err != nil {
        log.Fatal(err)
    }

    // Ping the database to verify connection
    err = client.Ping(context.TODO(), nil)
    if err != nil {
        log.Fatal(err)
    }

    log.Println("Connected to MongoDB!")
}
