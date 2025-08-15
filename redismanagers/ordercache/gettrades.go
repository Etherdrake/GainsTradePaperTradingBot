package ordercache

import (
	"HootTelegram/api"
	"github.com/go-redis/redis/v8"
	"log"
	"strings"
)

func GetAllPositionsRdb(rdb *redis.Client) ([]api.OpenTradeJSON, error) {
	pattern := "user:*:trade:*" // Adjust the pattern to match all trade keys

	var cursor uint64
	var trades []api.OpenTradeJSON

	for {
		var keys []string
		var err error
		keys, cursor, err = rdb.Scan(Ctx, cursor, pattern, 10).Result()
		if err != nil {
			return nil, err
		}
		for _, key := range keys {
			// Extract userGUID from the key
			keyParts := strings.Split(key, ":")
			if len(keyParts) < 4 {
				continue // Skip invalid keys
			}
			guid := keyParts[1]
			tradeId := keyParts[3]

			trade, getTradeErr := GetTradeSplit(rdb, guid, tradeId)
			if getTradeErr != nil {
				log.Println("getTradeFromKey Error:", getTradeErr)
				return nil, getTradeErr
			}
			trades = append(trades, trade)
		}
		if cursor == 0 {
			break
		}
	}
	return trades, nil
}

// GetAllOrdersForUserArrayFromRedis retrieves all trades / orders / positions
// from the Redis papertrade instance.
func GetAllOrdersForUserArrayFromRedis(rdbPositionsPaper *redis.Client, guid string) ([]api.OpenTradeJSON, error) {
	pattern := "user:" + guid + ":trade:*"

	var cursor uint64
	var trades []api.OpenTradeJSON

	for {
		var keys []string
		var err error
		keys, cursor, err = rdbPositionsPaper.Scan(Ctx, cursor, pattern, 10).Result()
		if err != nil {
			return nil, err
		}
		for _, key := range keys {
			trade, err := GetTradeWithHash(rdbPositionsPaper, guid, key)
			if err != nil {
				log.Println("getTradeFromKey Error:", err)
				return nil, err
			}
			trades = append(trades, trade)
		}
		if cursor == 0 {
			break
		}
	}
	return trades, nil
}
