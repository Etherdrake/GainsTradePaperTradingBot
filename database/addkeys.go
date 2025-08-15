package database

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func AddPrivateKey(userID int64, key string) error {
	// Create a MongoDB client
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		return err
	}

	// Disconnect from MongoDB when done
	defer client.Disconnect(context.Background())

	// Access the "users" collection in the "strikerdb" database
	collection := client.Database(databaseName).Collection(collectionName)

	// Define the update filter to find the user by userID
	filter := bson.M{"_id": userID}

	// Define the update fields
	update := bson.M{
		"$set": bson.M{
			fmt.Sprintf("wallet_set"):  true,
			fmt.Sprintf("private_key"): key,
		},
	}

	// Perform the update operation
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

func AddPublicKey(userID int64, key string) error {
	// Create a MongoDB client
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		return err
	}

	// Disconnect from MongoDB when done
	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(ctx)
		if err != nil {
			log.Println("Error disconnecting client in AddPublicKey: ", err)
		}
	}(client, context.Background())

	// Access the "users" collection in the "strikerdb" database
	collection := client.Database(databaseName).Collection(collectionName)

	// Define the update filter to find the user by userID
	filter := bson.M{"_id": userID}

	// Define the update fields
	update := bson.M{
		"$set": bson.M{
			fmt.Sprintf("wallet_set"): true,
			fmt.Sprintf("public_key"): key,
		},
	}

	// Perform the update operation
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return nil
}
