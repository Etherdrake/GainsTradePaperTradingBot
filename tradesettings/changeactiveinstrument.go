package tradesettings

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func IncrementActiveInstrument(userID int64) error {
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		return fmt.Errorf("failed to create MongoDB client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Printf("Failed to disconnect from MongoDB: %v", err)
		}
	}()

	collection := client.Database(databaseName).Collection(collectionName)

	filter := bson.D{{"userid", userID}}
	update := bson.D{{"$inc", bson.D{{"active_instrument", 1}}}}
	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to increment active instrument: %v", err)
	}

	return nil
}

// DecrementActiveInstrument decrements the active instrument count by one for the given userID in the MongoDB collection.
func DecrementActiveInstrument(userID int64) error {
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		return fmt.Errorf("failed to create MongoDB client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Printf("Failed to disconnect from MongoDB: %v", err)
		}
	}()

	collection := client.Database(databaseName).Collection(collectionName)

	filter := bson.D{{"userid", userID}}
	update := bson.D{{"$inc", bson.D{{"active_instrument", -1}}}}
	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to decrement active instrument: %v", err)
	}

	return nil
}
