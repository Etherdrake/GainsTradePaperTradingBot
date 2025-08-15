package ordercache

import (
	"HootTelegram/api"
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

var Ctx = context.Background()

func ConnectToRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	_, err := rdb.Ping(Ctx).Result()
	if err != nil {
		panic(err)
	}
	return rdb
}

// StoreOpenTradeInRedis stores a trade from the cache trade into Redis after it has been converted into api.OpenTradeJSON prior to
// Execution of this function
func StoreOpenTradeInRedis(rdb *redis.Client, userGUID string, trade api.OpenTradeJSON) error {
	key := "user:" + userGUID + ":trade:" + trade.TradeID

	// Get the current UNIX timestamp
	currentTimestamp := time.Now().Unix()

	fields := map[string]interface{}{
		"ID":              userGUID,
		"TradeID":         trade.TradeID,
		"Trader":          trade.Trader,
		"PairIndex":       trade.PairIndex,
		"Index":           trade.Index,
		"InitialPosToken": trade.InitialPosToken,
		"PositionSizeDai": trade.PositionSizeDai,
		"OpenPrice":       trade.OpenPrice,
		"Buy":             strconv.FormatBool(trade.Buy),
		"Leverage":        trade.Leverage,
		"TP":              trade.TP,
		"SL":              trade.SL,
		"PercentageTP":    trade.PercentageTP,
		"PercentageSL":    trade.PercentageSL,
		"OrderType":       trade.OrderType,
		"OrderStatus":     trade.OrderStatus,
		"Timestamp":       currentTimestamp, // Add this line
	}

	return rdb.HMSet(Ctx, key, fields).Err()
}

func ChangeOrderStatus(rdb *redis.Client, userGUID string, tradeID string, newOrderStatus string) error {
	key := "user:" + userGUID + ":trade:" + tradeID

	// Update the OrderStatus field
	fields := map[string]interface{}{
		"OrderStatus": newOrderStatus,
	}

	return rdb.HMSet(Ctx, key, fields).Err()
}

func PopTradeFromRedis(rdb *redis.Client, userGUID string, tradeID string) (api.OpenTradeJSON, error) {
	// Firstly, get the trade using GetOpenTradeFromRedis
	trade, err := GetOpenTradeFromRedis(rdb, userGUID, tradeID)
	if err != nil {
		return api.OpenTradeJSON{}, err
	}

	// Now, delete the trade from Redis
	key := "user:" + userGUID + ":trade:" + tradeID
	_, err = rdb.Del(Ctx, key).Result()
	if err != nil {
		return api.OpenTradeJSON{}, err
	}

	// Return the trade that was removed
	return trade, nil
}

func GetOpenTradeFromRedis(rdbPos *redis.Client, userGUID string, tradeID string) (api.OpenTradeJSON, error) {
	key := "user:" + userGUID + ":trade:" + tradeID

	fields, err := rdbPos.HGetAll(Ctx, key).Result()
	if err != nil {
		return api.OpenTradeJSON{}, err
	}

	// Convert Buy field to bool
	buyBool, err := strconv.ParseBool(fields["Buy"])
	if err != nil {
		return api.OpenTradeJSON{}, err
	}

	// Convert Ordertype field to uint8
	ordertypeStr, ok := fields["OrderType"]
	if !ok {
		return api.OpenTradeJSON{}, errors.New("OrderType field not found in Redis hash")
	}

	ordertypeValue, err := strconv.ParseUint(ordertypeStr, 10, 8)
	if err != nil {
		return api.OpenTradeJSON{}, err
	}

	orderTimestamp, err := strconv.ParseUint(fields["Timestamp"], 10, 64)
	if err != nil {
		return api.OpenTradeJSON{}, err
	}

	orderStatus, err := strconv.ParseUint(fields["OrderStatus"], 10, 64)
	if err != nil {
		return api.OpenTradeJSON{}, err
	}

	userGuid, err := strconv.ParseInt(userGUID, 10, 64)
	if err != nil {
		return api.OpenTradeJSON{}, err
	}

	return api.OpenTradeJSON{
		ID:              userGuid,
		TradeID:         fields["TradeID"],
		Trader:          fields["Trader"],
		PairIndex:       fields["PairIndex"],
		Index:           fields["Index"],
		InitialPosToken: fields["InitialPosToken"],
		PositionSizeDai: fields["PositionSizeDai"],
		OpenPrice:       fields["OpenPrice"],
		Buy:             buyBool,
		Leverage:        fields["Leverage"],
		TP:              fields["TP"],
		SL:              fields["SL"],
		OrderType:       uint8(ordertypeValue),
		OrderStatus:     uint8(orderStatus),
		Timestamp:       orderTimestamp,
	}, nil
}
