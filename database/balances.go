package database

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func AddPaperPnl(userID int64, amount float64) error {
	return updateBalanceOrPnl(userID, "paper_pnl", amount)
}

func AddPaperBalance(userID int64, amount float64) error {
	return updateBalanceOrPnl(userID, "paper_balance", amount)
}

func AddPerpPnl(userID int64, amount float64) error {
	return updateBalanceOrPnl(userID, "perp_pnl", amount)
}

func AddPerpBalance(userID int64, amount float64) error {
	return updateBalanceOrPnl(userID, "perp_balance", amount)
}

func AddXP(userID int64, amount float64) error {
	return updateBalanceOrPnl(userID, "total_xp", amount)
}

func DecrXP(userID int64, amount float64) error {
	return updateBalanceOrPnl(userID, "total_xp", -amount)
}

func DecrPaperPnl(userID int64, amount float64) error {
	return updateBalanceOrPnl(userID, "paper_pnl", -amount)
}

func DecrPaperBalance(userID int64, amount float64) error {
	return updateBalanceOrPnl(userID, "paper_balance", amount)
}

func DecrPerpPnl(userID int64, amount float64) error {
	return updateBalanceOrPnl(userID, "perp_pnl", -amount)
}

func DecrPerpBalance(userID int64, amount float64) error {
	return updateBalanceOrPnl(userID, "perp_balance", -amount)
}

func GetPaperPnl(userID int64) (float64, error) {
	return GetFieldForUser(userID, "paper_pnl")
}

func GetPaperBalance(userID int64) (float64, error) {
	return GetFieldForUser(userID, "paper_balance")
}

func GetPerpPnl(userID int64) (float64, error) {
	return GetFieldForUser(userID, "perp_pnl")
}

func GetPerpBalance(userID int64) (float64, error) {
	return GetFieldForUser(userID, "perp_balance")
}

func GetFieldForUser(userID int64, field string) (float64, error) {
	// Create a MongoDB client
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		return 0, err
	}
	defer client.Disconnect(context.Background())

	// Access the "users" collection in the "strikerdb" database
	collection := client.Database(databaseName).Collection(collectionName)

	// Define the filter to find the user by userID
	filter := bson.M{
		"_id": userID,
	}

	// Find the user document
	var user User
	err = collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return 0, err
	}

	// Return the requested field value
	switch field {
	case "paper_pnl":
		return user.PaperPnl, nil
	case "paper_balance":
		return user.PaperBalance, nil
	case "perp_pnl":
		return user.PerpPnl, nil
	case "perp_balance":
		return user.PerpBalance, nil
	default:
		return 0, fmt.Errorf("invalid field name")
	}
}

func updateBalanceOrPnl(userID int64, field string, amount float64) error {
	// Create a MongoDB client
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		return err
	}
	defer client.Disconnect(context.Background())

	// Access the "users" collection in the "strikerdb" database
	collection := client.Database(databaseName).Collection(collectionName)

	// Define the filter to find the user by userID
	filter := bson.M{
		"_id": userID,
	}

	//fmt.Println("Mongo Update Amount: ", amount)

	// Define the update operation to increment/decrement the specified field
	update := bson.M{
		"$inc": bson.M{
			field: amount,
		},
	}

	// Perform the update operation
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return nil
}
