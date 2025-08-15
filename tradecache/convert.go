package tradecache

import (
	"HootTelegram/api"
	"strconv"
)

func ConvertToOpenTradeJSON(cacheTrade OpenTradeCache) api.OpenTradeJSON {
	return api.OpenTradeJSON{
		ID:              cacheTrade.ID,
		TradeID:         cacheTrade.TradeID,
		Trader:          cacheTrade.Trader,
		PairIndex:       strconv.FormatInt(cacheTrade.PairIndex, 10),
		Index:           strconv.FormatInt(cacheTrade.Index, 10),
		InitialPosToken: strconv.FormatInt(cacheTrade.InitialPosToken, 10),
		PositionSizeDai: strconv.FormatInt(cacheTrade.PositionSizeDai, 10),
		OpenPrice:       strconv.FormatFloat(cacheTrade.OpenPrice, 'f', 10, 64),
		ClosePrice:      strconv.FormatFloat(cacheTrade.ClosePrice, 'f', 10, 64), // assuming OpenPrice is a float64, and you want 2 decimal places
		Buy:             cacheTrade.Buy,
		Leverage:        strconv.FormatInt(int64(cacheTrade.Leverage), 10),
		TP:              strconv.FormatFloat(cacheTrade.TP, 'f', 10, 64), // assuming TP is a float64, and you want 2 decimal places
		SL:              strconv.FormatFloat(cacheTrade.SL, 'f', 10, 64), // assuming SL is a float64, and you want 2 decimal places
		PercentageTP:    strconv.FormatFloat(cacheTrade.PercentageTP, 'f', 2, 64),
		PercentageSL:    strconv.FormatFloat(cacheTrade.PercentageSL, 'f', 2, 64),
		OrderType:       cacheTrade.OrderType,   // Replace with the actual OrderType value
		OrderStatus:     cacheTrade.OrderStatus, // Replace with the actual OrderStatus value
		Chain:           "Practice",
	}
}
