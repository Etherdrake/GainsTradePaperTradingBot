package ordercache

import (
	"HootTelegram/api"
	"errors"
	"github.com/go-redis/redis/v8"
	"strconv"
)

func GetTradeSplit(rdb *redis.Client, guid string, orderId string) (api.OpenTradeJSON, error) {
	fields, err := rdb.HGetAll(ctx, "user:"+guid+":trade:"+orderId).Result()
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
		return api.OpenTradeJSON{}, errors.New("ordertype field not found in Redis hash")
	}

	ordertypeValue, err := strconv.ParseUint(ordertypeStr, 10, 8)
	if err != nil {
		return api.OpenTradeJSON{}, err
	}

	ID, err := strconv.ParseInt(guid, 10, 64)
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

	var orderPnL float64
	if fields["PnL"] != "" {
		orderPnL, err = strconv.ParseFloat(fields["PnL"], 64)
		if err != nil {
			return api.OpenTradeJSON{}, err
		}
	} else {
		orderPnL = 0
	}

	return api.OpenTradeJSON{
		ID:              ID,
		TradeID:         orderId,
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
		PnL:             orderPnL,
		Chain:           fields["Chain"],
		Timestamp:       orderTimestamp,
	}, nil
}
