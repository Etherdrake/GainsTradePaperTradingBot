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

// GetLongShort retrieves `true` if "long" is true and `false` if "long" is false for the given userID from the MongoDB collection.
func GetLongShort(userID int64) (bool, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		return false, fmt.Errorf("failed to create MongoDB client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to connect to MongoDB: %v", err)
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
			return false, fmt.Errorf("no trade settings found for the given userID")
		}
		return false, fmt.Errorf("failed to find trade settings: %v", err)
	}

	return tradeSetting.Long, nil
}
