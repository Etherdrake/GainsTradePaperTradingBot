package main

import (
	"HootTelegram/redismanagers/ordercache"
	"HootTelegram/tradecache"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
)

// SetStartTradeOrOrder sets the first trade or order as active in the cache. TODO
func SetStartTradeOrOrder(client *mongo.Client, tradeCache *tradecache.TradeCache, guid int64, isArb bool) error {
	return nil
}

func SetTradeId(tradeCache *tradecache.TradeCache, rdbPositionsPaper *redis.Client, guid int64) error {
	tradeArr, err := ordercache.GetAllPositionsRdb(rdbPositionsPaper)
	if err != nil {
		return err
	}
	if len(tradeArr) > 0 {
		err := tradeCache.SetActiveTradeID(guid, tradeArr[0].TradeID)
		if err != nil {
			return err
		}
	}
	return nil
}
