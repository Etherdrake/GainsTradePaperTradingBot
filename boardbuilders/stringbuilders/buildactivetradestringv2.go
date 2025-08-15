package stringbuilders

import (
	"HootTelegram/redismanagers/ordercache"
	"HootTelegram/tradecache"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"strconv"
)

func BuildActiveGainsTradeStringV2(
	tradeCache *tradecache.TradeCache,
	rdbPrice *redis.Client, rdbPositionPaper *redis.Client,
	activeTradeID string,
	guid int64) string {

	trade, err := ordercache.GetTradeSplit(rdbPrice, strconv.FormatInt(guid, 10), activeTradeID)
	if err != nil {
		// LOGGING_ALERT
		log.Println("No trade was found ")
	}

	var openTradesMsg string
	var openOrdersMsg string
	var totalTrades int
	var totalOrders int

	// This builds all the strings you want EXCEPT for the ActiveTradeString which you ALWAYS want
	if ordercache.HasTrades(rdbPositionPaper, strconv.FormatInt(guid, 10)) {
		openTradesMsg, openOrdersMsg, totalTrades, totalOrders, err = BuildPositionsAndOrderString(rdbPrice, rdbPositionPaper, activeTradeID, guid)
		if err != nil {
			log.Println("Error in buildpositionsandorderstring in buildperpmainstring.go for user", guid)
		}
	} else {
		openTradesMsg = "You don't have any open trades."
	}

	//pairStr, _ := strconv.Atoi(trade.PairIndex)
	//
	//// Retrieve the price from Redis
	//price, err := priceserver.GetPrice(rdbPrice, pairStr) // Assuming you have the GetPrice function implemented to fetch price from Redis
	//if err != nil {
	//	fmt.Println("Error fetching price from Redis in Stringbuilder:", err)
	//}

	//var LongShortString string
	//if trade.Buy {
	//	LongShortString = "Long"
	//} else {
	//	LongShortString = "Short"
	//}

	totalActive := totalTrades

	if totalTrades != 0 {
		totalActive += 1
	}

	var activeTradeOrOrder string
	// We retrieve the activeTradeOrOrder
	activeTradeOrOrder = BuildOrderOrTradeStringPaperV3(rdbPrice, trade)

	var msg string
	// We have no trades but an order
	if totalTrades == 0 && totalOrders > 0 {
		msg = fmt.Sprintf(`
⭐ *All Active Trades (%d)*:

%s
%s

%s`,
			totalActive,
			activeTradeOrOrder,
			openOrdersMsg,
			openTradesMsg,
		)
	} else {
		// We have trades and an order
		if totalTrades > 0 && totalOrders > 0 {
			msg = fmt.Sprintf(`
⭐ *All Active Trades (%d)*:
%s

%s
%s`,
				totalActive,
				activeTradeOrOrder,
				openTradesMsg,
				openOrdersMsg,
			)
		} else {
			// We have trades and no order
			if totalTrades > 0 && totalOrders == 0 {
				msg = fmt.Sprintf(`
⭐ *All Active Trades (%d)*:
%s

%s`,
					totalActive,
					activeTradeOrOrder,
					openTradesMsg,
				)
			}
		}
	}
	return msg
}
