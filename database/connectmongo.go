package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectToMongoDB(uri, dbName string) (*mongo.Database, error) {
	// Set up MongoDB connection options
	clientOptions := options.Client().ApplyURI(uri)

	// Create a MongoDB client
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return nil, err
	}

	// Set up a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to the MongoDB server
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	// Return a handle to the specified database
	return client.Database(dbName), nil
}
