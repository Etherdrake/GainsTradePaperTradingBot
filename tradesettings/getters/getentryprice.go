package getters

import (
	"HootTelegram/types"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

// GetEntryPrice retrieves the "entry_price" value for the given userID from the MongoDB collection.
func GetEntryPrice(userID int64) (float64, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		return 0, fmt.Errorf("failed to create MongoDB client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Printf("Failed to disconnect from MongoDB: %v", err)
		}
	}()

	collection := client.Database(databaseName).Collection(collectionName)

	filter := bson.D{{"_id", userID}}
	var tradeSetting types.TradeSetting
	err = collection.FindOne(ctx, filter).Decode(&tradeSetting)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return 0, fmt.Errorf("no trade settings found for the given userID")
		}
		return 0, fmt.Errorf("failed to find trade settings: %v", err)
	}

	return tradeSetting.EntryPrice, nil
}
