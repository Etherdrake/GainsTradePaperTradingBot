package database

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func AddSettingsDB(userID int64, chain string) error {
	// Create a MongoDB client
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		return err
	}

	// Disconnect from MongoDB when done
	defer client.Disconnect(context.Background())

	// Access the "users" collection in the "strikerdb" database
	collection := client.Database(databaseName).Collection(chain + "_settings")

	// Create the document to insert
	document := bson.M{
		"_id": userID,

		// General
		"anti_rug":        false,
		"semi_private_tx": false,
		"smarts_lippage":  false,
		"multi_wallet":    false,
		"max_gas_price":   10,
		"slippage":        10,
		"max_gas_limit":   3_000_000,
		"degen_mode":      false,

		// BUY
		"auto_buy":      false,
		"duplicate_buy": false,
		"buy_gas_price": 3,
		"max_mcap":      false,
		"min_liq":       false,
		"max_liq":       false,
		"min_mcap":      false,
		"max_buy_tax":   false,
		"max_sell_tax":  false,

		// SELL
		"auto_sell":               false,
		"trading_sell":            false,
		"trade_sell_confirmation": false,
		"sell_gas_price":          3,
		"auto_sell_high":          100,
		"sell_amount_high":        100,
		"auto_sell_low":           -101,
		"sell_amount_low":         100,

		// APPROVE
		"auto_approve": false,
		// Default gas price is 3 gwei
		"approve_gas_price": 3,

		// Striker Fees
		"unpaid": 0,
	}

	// Insert the document into the collection
	_, err = collection.InsertOne(context.Background(), document)
	if err != nil {
		return err
	}

	return nil
}
