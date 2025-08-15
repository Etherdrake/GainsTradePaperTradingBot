package ordercache

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
)

var ctx = context.Background()

func HasTrades(rdb *redis.Client, guid string) bool {
	pattern := fmt.Sprintf("user:%s:trade:*", guid)

	// Using the SCAN command to avoid blocking the database with KEYS in large datasets.
	iter := rdb.Scan(ctx, 0, pattern, 0).Iterator()
	if iter.Next(ctx) {
		return true
	}

	if err := iter.Err(); err != nil {
		log.Printf("Error checking trades for user: %v", err)
		return false
	}

	return false
}
