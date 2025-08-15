package database

import (
	"HootTelegram/types"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func GetKeys(userID int64) (*types.Keys, error) {
	// Connect to MongoDB and access the 'users' collection in 'hooterdb' database
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return nil, err
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(ctx)

	// Get the 'users' collection
	users := client.Database("hooterdb").Collection("users")

	// Find the user with the given guid
	var result struct {
		PublicKey  string `bson:"public_key"`
		PrivateKey string `bson:"private_key"`
	}
	filter := bson.D{{"_id", userID}}
	err = users.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &types.Keys{PublicKey: result.PublicKey, PrivateKey: result.PrivateKey}, nil
}
