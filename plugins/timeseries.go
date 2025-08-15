package plugins

import (
	"HootTelegram/pairmaps"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var TimeFrames = []struct {
	Interval   time.Duration
	Collection string
}{
	{5 * time.Minute, "5M"},
	{15 * time.Minute, "15M"},
}

func GenerateTimeSeriesStorageAgent(rdbPrice *redis.Client, database *mongo.Database) error {
	ctx := context.Background()

	mongoOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	mongoClient, err := mongo.Connect(ctx, mongoOptions)
	if err != nil {
		return err
	}
	defer mongoClient.Disconnect(ctx)

	// Calculate the next start time with a 5-minute increment
	now := time.Now()
	nextStartTime := now.Add(time.Duration((5-now.Minute()%5)%5) * time.Minute)
	if now.After(nextStartTime) {
		nextStartTime = nextStartTime.Add(5 * time.Minute)
	}

	// Wait until the next start time
	timeUntilStart := nextStartTime.Sub(now)
	time.Sleep(timeUntilStart)

	// Loop through each currency pair and each timeframe
	for index, pair := range pairmaps.IndexToPair {
		for _, tf := range TimeFrames {
			go func(index int, pair string, interval time.Duration, collection string) {
				for {
					priceKey := fmt.Sprintf("price:%d", index)
					priceData, err := rdbPrice.Get(ctx, priceKey).Result()
					if err != nil {
						fmt.Printf("Error retrieving price data for %s: %s\n", pair, err)
						continue
					}

					// Store the price data in the appropriate MongoDB collection
					err = storePriceData(ctx, database, pair, priceData, collection)
					if err != nil {
						fmt.Printf("Error storing price data for %s (%s): %s\n", pair, collection, err)
					}

					// Wait for the next interval
					time.Sleep(interval)
				}
			}(index, pair, tf.Interval, tf.Collection)
		}
	}

	select {}
}

func storePriceData(ctx context.Context, database *mongo.Database, pair, priceData, interval string) error {
	collection := database.Collection(interval)
	document := bson.M{
		"pair":      pair,
		"price":     priceData,
		"timestamp": time.Now().UTC(), // Use UTC timestamp
	}

	_, err := collection.InsertOne(ctx, document)
	return err
}
