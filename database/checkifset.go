package database

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CheckIfSet(userID int64) (bool, error) {
	// Create a MongoDB client
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		return false, err
	}

	// Disconnect from MongoDB when done
	defer client.Disconnect(context.Background())

	// Access the "users" collection in the "strikerdb" database
	collection := client.Database(databaseName).Collection(collectionName)

	// Define the filter to find the user by userID and check if {chain}Set is true
	filter := bson.M{
		"_id":                     userID,
		fmt.Sprintf("wallet_set"): true,
	}

	// Count the matching documents in the collection
	count, err := collection.CountDocuments(context.Background(), filter)
	if err != nil {
		return false, err
	}

	// If count is greater than 0, it means the {chain}Set is true for the user
	isSet := count > 0

	return isSet, nil
}
