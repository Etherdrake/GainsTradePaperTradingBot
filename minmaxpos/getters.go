package minmaxpos

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"strconv"
	"strings"
)

func GetMinLeverage(index int, client *redis.Client) (int64, error) {
	result, err := client.Get(context.Background(), strconv.Itoa(index)).Result()
	if err != nil {
		return 0, err
	}

	parts := strings.Split(result, ";")
	if len(parts) < 3 {
		return 0, errors.New("invalid data format in Redis")
	}

	minimum, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return 0, err
	}

	return minimum, nil
}

func GetMaxLeverage(index int, client *redis.Client) (int64, error) {
	result, err := client.Get(context.Background(), strconv.Itoa(index)).Result()
	if err != nil {
		return 0, err
	}

	parts := strings.Split(result, ";")
	if len(parts) < 3 {
		return 0, errors.New("invalid data format in Redis")
	}

	maximum, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return 0, err
	}

	return maximum, nil
}

func GetMinPositionSize(index int, client *redis.Client) (int64, error) {
	result, err := client.Get(context.Background(), strconv.Itoa(index)).Result()
	if err != nil {
		return 0, err
	}

	parts := strings.Split(result, ";")
	if len(parts) < 3 {
		return 0, errors.New("invalid data format in Redis")
	}

	minLevPosDai, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return 0, err
	}

	return minLevPosDai, nil
}
