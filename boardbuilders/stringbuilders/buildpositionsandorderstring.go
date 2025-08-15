package stringbuilders

import (
	"HootTelegram/redismanagers/ordercache"
	"github.com/go-redis/redis/v8"
	"log"
	"strconv"
)

// BuildPositionsAndOrderString build a string that shows: All Active Trades and All Open Orders.
// We use the following logic:
// 1. Get all positions and orders for user
// 2. Remove all non-open trades and orders
// 3. Seperate open_trades and open_orders
// 4. Build the corresponding string using the arrays as input
func BuildPositionsAndOrderString(
	rdbPrice *redis.Client,
	rdbPaper *redis.Client,
	activeTradeID string,
	guid int64) (string, string, int, int, error) {
	var openTradesMsg string
	var openOrdersMsg string
	var totalTrades int
	var totalOrders int
	openTradesMsg = ""
	openOrdersMsg = ""

	// Check if user has open positions in Redis
	if ordercache.HasTrades(rdbPaper, strconv.FormatInt(guid, 10)) {
		allOrders, err := ordercache.GetAllOrdersForUserArrayFromRedis(rdbPaper, strconv.FormatInt(guid, 10))
		if err != nil {
			log.Println("Error encountered get AllOrdersArray in BuildPositionsAndOrderString", err)
		}
		// Seperate open_trades and open_orders
		openTradesArray := ordercache.FilterOpenTrades(allOrders)
		openOrdersArray := ordercache.FilterOpenOrders(allOrders)

		// Sort both arrays
		openTradesArraySorted := ordercache.SortArrayByTimestamp(openTradesArray)
		openOrdersArraySorted := ordercache.SortArrayByTimestamp(openOrdersArray)

		// Build the two strings, but leave the activeID out:
		for i := 0; i < len(openTradesArraySorted); i++ {
			if openTradesArraySorted[i].TradeID == activeTradeID {
				continue
			}
			openTradesStr, err := BuildSinglePositionString(rdbPrice, openTradesArraySorted[i], i)
			if err != nil {
				log.Println("Error getting generating position string in BuildPositionsAndOrderString")
			}
			if i != len(openTradesArraySorted) && i != len(openTradesArraySorted) {
				openTradesMsg += openTradesStr + "\n"
			} else {
				openTradesMsg += openTradesStr
			}

			totalTrades += 1
		}
		for x := 0; x < len(openOrdersArraySorted); x++ {
			if openTradesArraySorted[x].TradeID == activeTradeID {
				continue
			}
			openOrderStr, err := BuildSingleOrderString(rdbPrice, openOrdersArraySorted[x], x)
			if err != nil {
				log.Println("Error getting generating order string in BuildPositionsAndOrderString")
			}
			if x != len(openTradesArraySorted) && x != len(openTradesArraySorted) {
				openOrdersMsg += openOrderStr + "\n"
			} else {
				openOrdersMsg += openOrderStr
			}

			totalOrders += 1
		}
	}
	return openTradesMsg, openOrdersMsg, totalTrades, totalOrders, nil
}
