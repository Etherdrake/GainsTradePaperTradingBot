package redislocker

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"strconv"
)

// CheckUserID checks if a user ID exists in the active users set in Redis
func CheckUserID(ctx context.Context, client *redis.Client, userID int64) (bool, error) {
	exists, err := client.SIsMember(ctx, "active_users", strconv.FormatInt(userID, 10)).Result()
	if err != nil {
		log.Printf("Error checking if user ID %d exists in active users set: %v", userID, err)
		return false, err
	}
	return exists, nil
}
