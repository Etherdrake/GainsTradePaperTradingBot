package database

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitUser(userID int64, userName, firstName, lastName string) error {
	// Create a MongoDB client
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		return err
	}

	// Disconnect from MongoDB when done
	defer client.Disconnect(context.Background())

	// Access the "users" collection in the "strikerdb" database
	collection := client.Database(databaseName).Collection(collectionName)

	// Create the document to insert
	document := bson.M{
		"_id":           userID,
		"user_name":     userName,
		"first_name":    firstName,
		"last_name":     lastName,
		"wallet_set":    false,
		"public_key":    "",
		"private_key":   "",
		"paper_pnl":     0,
		"paper_balance": 10000,
		"perp_pnl":      0,
		"perp_balance":  0,
		"total_xp":      0,
	}

	// Insert the document into the collection
	_, err = collection.InsertOne(context.Background(), document)
	if err != nil {
		return err
	}

	return nil
}
