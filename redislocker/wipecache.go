package redislocker

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
)

// WipeCache removes all keys from Redis
func WipeCache(ctx context.Context, client *redis.Client) error {
	// FLUSHDB command removes all keys from the currently selected Redis database
	err := client.FlushDB(ctx).Err()
	if err != nil {
		log.Printf("Error wiping cache: %v", err)
		return err
	}
	return nil
}
