package tradesettings

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func InitTradeSettings(userID int64) error {
	// Create a MongoDB client
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		return err
	}

	// Disconnect from MongoDB when done
	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(ctx)
		if err != nil {
			log.Println("Error in mongo Disconnect inside InitTradeSettings: ", err)
		}
	}(client, context.Background())

	// Access the "tradesettings" collection in the "your_database_name" database
	collection := client.Database(databaseName).Collection(collectionName)

	// Check if the document already exists for the given userID
	filter := bson.D{{"_id", userID}}
	count, err := collection.CountDocuments(context.Background(), filter, nil)
	if err != nil {
		return err
	}

	// If the document already exists, return without inserting a new one
	if count > 0 {
		return nil
	}

	// Create the document to insert
	document := bson.M{
		"_id":               userID,
		"entry_price":       25000,
		"active_instrument": 0,
		"position_size":     100,
		"long":              true,
		"leverage":          5, // Initialize with the default leverage value (e.g., 1)
		"take_profit":       0, // Initialize with the default take profit value (e.g., 0)
		"stop_loss":         0, // Initialize with the default stop loss value (e.g., 0)
	}

	// Insert the document into the collection
	_, err = collection.InsertOne(context.Background(), document)
	if err != nil {
		return err
	}

	return nil
}
