package database

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
	"time"
)

func GetPublicKey(userID int64) (string, error) {
	// Connect to MongoDB and access the 'users' collection in 'hooterdb' database

	// Replace "mongodb://localhost:27017" with your MongoDB connection string
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return "", err
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		return "", err
	}
	defer client.Disconnect(ctx)

	// Get the 'users' collection
	users := client.Database("hooterdb").Collection("users")

	// Find the user with the given guid
	var result struct {
		PublicKey string `bson:"public_key"`
	}
	filter := bson.D{{"_id", userID}}
	err = users.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return "", err
	}

	return strings.ToLower(result.PublicKey), nil
}
