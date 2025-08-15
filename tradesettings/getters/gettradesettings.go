package getters

import (
	"HootTelegram/types"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

// GetTradeSettingsByUserID takes a userID int64 and returns the trade settings document associated with it.
func GetTradeSettingsByUserID(userID int64) (*types.TradeSetting, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, fmt.Errorf("failed to create MongoDB client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
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
		if err == mongo.ErrNoDocuments {
			return nil, nil // No document found for the given userID
		}
		return nil, fmt.Errorf("failed to find trade settings: %v", err)
	}

	return &tradeSetting, nil
}
