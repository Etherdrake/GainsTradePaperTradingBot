package openpositionedit

import (
	"HootTelegram/api"
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
)

// SetZeroTP sets the TP to zero for an OpenTradeJSON object.
func SetZeroTP(trade api.OpenTradeJSON) api.OpenTradeJSON {
	trade.TP = "0"
	return trade
}

// DecrPosTP decreases the TP by 0.25% while taking into account the leverage.
func DecrPosTP(trade api.OpenTradeJSON) (api.OpenTradeJSON, error) {
	// Convert TP to float64 for calculation
	tp, err := strconv.ParseFloat(trade.TP, 64)
	if err != nil {
		return trade, err
	}

	// Calculate the decrease amount
	decrease := tp * 0.0025

	// Ensure TP doesn't go negative
	if tp-decrease < 0 {
		return trade, errors.New("TP cannot go negative")
	}

	// Update TP
	tp -= decrease

	// Convert back to string and update trade
	trade.TP = fmt.Sprintf("%.2f", tp)
	return trade, nil
}

// IncrPosTP increases the TP by 0.25% while taking into account the leverage.
func IncrPosTP(trade api.OpenTradeJSON) (api.OpenTradeJSON, error) {
	// Convert TP to float64 for calculation
	tp, err := strconv.ParseFloat(trade.TP, 64)
	if err != nil {
		return trade, err
	}

	// Calculate the increase amount
	increase := tp * 0.0025

	// Update TP
	tp += increase

	// Convert back to string and update trade
	trade.TP = fmt.Sprintf("%.2f", tp)
	return trade, nil
}

// UpdateTPInCache updates the TP value of a trade in the Redis cache.
func UpdateTPInCache(rdbPos *redis.Client, guid string, tradeID string, newTp string) error {
	// Assuming your Redis key format is the same as in your original function
	key := "user:" + guid + ":trade:" + tradeID

	// Update the TP value in the Redis hash
	_, err := rdbPos.HSet(context.TODO(), key, "TP", newTp).Result()
	if err != nil {
		return err
	}
	return nil
}
