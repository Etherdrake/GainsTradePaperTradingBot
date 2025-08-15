package openpositionedit

import (
	"HootTelegram/api"
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
)

// SetZeroSL sets the SL to zero for an OpenTradeJSON object.
func SetZeroSL(trade api.OpenTradeJSON) api.OpenTradeJSON {
	trade.SL = "0"
	return trade
}

// DecrPosSL decreases the SL by 0.25%.
func DecrPosSL(trade api.OpenTradeJSON) (api.OpenTradeJSON, error) {
	// Convert SL to float64 for calculation
	sl, err := strconv.ParseFloat(trade.SL, 64)
	if err != nil {
		return trade, err
	}

	// Calculate the decrease amount
	decrease := sl * 0.0025

	// Ensure SL doesn't go negative
	if sl-decrease < 0 {
		return trade, errors.New("SL cannot go negative")
	}

	// Update SL
	sl -= decrease

	// Convert back to string and update trade
	trade.SL = fmt.Sprintf("%.2f", sl)
	return trade, nil
}

// IncrPosSL increases the SL by 0.25%.
func IncrPosSL(trade api.OpenTradeJSON) (api.OpenTradeJSON, error) {
	// Convert SL to float64 for calculation
	sl, err := strconv.ParseFloat(trade.SL, 64)
	if err != nil {
		return trade, err
	}

	// Calculate the increase amount
	increase := sl * 0.0025

	// Update SL
	sl += increase

	// Convert back to string and update trade
	trade.SL = fmt.Sprintf("%.2f", sl)
	return trade, nil
}

// UpdateSLInCache updates the SL value of a trade in the Redis cache.
func UpdateSLInCache(rdbPos *redis.Client, guid string, tradeID string, newSl string) error {
	// Assuming your Redis key format is the same as in your original function
	key := "user:" + guid + ":trade:" + tradeID

	// Update the SL value in the Redis hash
	_, err := rdbPos.HSet(context.TODO(), key, "SL", newSl).Result()
	if err != nil {
		return err
	}
	return nil
}
