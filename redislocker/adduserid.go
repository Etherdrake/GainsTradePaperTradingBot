package redislocker

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"strconv"
)

// AddUserID adds a user ID to the active users set in Redis
func AddUserID(ctx context.Context, client *redis.Client, userID int64) error {
	err := client.SAdd(ctx, "active_users", strconv.FormatInt(userID, 10)).Err()
	if err != nil {
		log.Printf("Error adding user ID %d to active users set: %v", userID, err)
		return err
	}
	return nil
}
