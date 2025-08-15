package ordercache

import (
	"HootTelegram/api"
	"github.com/go-redis/redis/v8"
	"strconv"
)

func SortArrayByTimestamp(trades []api.OpenTradeJSON) []api.OpenTradeJSON {
	quickSortByTimestamp(trades, 0, len(trades)-1)
	return trades
}

func SortPositionsByTimestamp(rdb *redis.Client, userGUID int64) ([]api.OpenTradeJSON, error) {
	guidStr := strconv.FormatInt(userGUID, 10)
	trades, err := GetAllOrdersForUserArrayFromRedis(rdb, guidStr)
	if err != nil {
		return nil, err
	}

	quickSortByTimestamp(trades, 0, len(trades)-1)

	return trades, nil
}

func quickSortByTimestamp(trades []api.OpenTradeJSON, low, high int) {
	if low < high {
		pivotIdx := partitionByTimestamp(trades, low, high)
		quickSortByTimestamp(trades, low, pivotIdx-1)
		quickSortByTimestamp(trades, pivotIdx+1, high)
	}
}

func partitionByTimestamp(trades []api.OpenTradeJSON, low, high int) int {
	pivot := trades[high].Timestamp
	i := low - 1

	for j := low; j < high; j++ {
		if trades[j].Timestamp >= pivot {
			i++
			trades[i], trades[j] = trades[j], trades[i]
		}
	}

	trades[i+1], trades[high] = trades[high], trades[i+1]
	return i + 1
}
