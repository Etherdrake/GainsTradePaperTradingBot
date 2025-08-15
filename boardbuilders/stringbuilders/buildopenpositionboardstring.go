package stringbuilders

import (
	"HootTelegram/api"
	"HootTelegram/pairmaps"
	"HootTelegram/priceserver"
	"HootTelegram/redismanagers/ordercache"
	"HootTelegram/utils"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"strconv"
)

func BuildOpenPositionBoardString(rdbPaper *redis.Client, rdbPrice *redis.Client, activeIdx int, guid int64) string {
	var activeTrade api.OpenTradeJSON
	var openTradesArraySorted []api.OpenTradeJSON

	// Check if this is about an openTrade or openOrder
	// Check if user has open positions in Redis
	if ordercache.HasTrades(rdbPaper, strconv.FormatInt(guid, 10)) {
		allOrders, err := ordercache.GetAllOrdersForUserArrayFromRedis(rdbPaper, strconv.FormatInt(guid, 10))
		if err != nil {
			log.Println("Error encountered get AllOrdersArray in BuildPositionsAndOrderString", err)
		}
		// Seperate open_trades and open_orders
		openTradesArray := ordercache.FilterOpenTrades(allOrders)

		// Sort the array
		openTradesArraySorted = ordercache.SortArrayByTimestamp(openTradesArray)

		// Get the active position using the idx provided
		activeTrade = openTradesArraySorted[activeIdx]

	}

	openPrice, err := strconv.ParseFloat(activeTrade.OpenPrice, 64)
	if err != nil {
		return ""
	}

	positionSize, err := strconv.ParseFloat(activeTrade.PositionSizeDai, 64)
	openPrice, err = strconv.ParseFloat(activeTrade.OpenPrice, 64)
	takeProfit, err := strconv.ParseFloat(activeTrade.TP, 64)
	pairIdx, _ := strconv.Atoi(activeTrade.PairIndex)
	leverage, err := strconv.ParseFloat(activeTrade.Leverage, 64)

	// Retrieve the price from Redis
	currentPrice, err := priceserver.GetPrice(rdbPrice, pairIdx) // Assuming you have the GetPrice function implemented to fetch price from Redis
	if err != nil {
		fmt.Println("Error fetching price from Redis in Stringbuilder:", err)
	}

	// Retrieve the ticker
	ticker := pairmaps.IndexToPair[pairIdx]

	getPnlUSD := utils.CalculateDollarProfitOrLossWithLeverage(positionSize, openPrice, currentPrice, leverage, activeTrade.Buy)
	getPnlPrcnt := utils.CalculatePercentageProfitOrLoss(positionSize, getPnlUSD, activeTrade.Buy)

	var getPnlOperator string
	if getPnlUSD > 0 {
		getPnlOperator = "+"
	}

	tradeDirection := "📉 Short"
	if activeTrade.Buy {
		tradeDirection = " 📈 Long"
	}

	liqPrice := utils.CalculateLiquidationPrice(openPrice, leverage, activeTrade.Buy)
	payOut := utils.CalculatePayout(openPrice, takeProfit, positionSize, leverage, activeTrade.Buy)
	if payOut < 0 {
		payOut = payOut * -1
	}

	activeIdxIncr := activeIdx + 1

	activePositionString := fmt.Sprintf("🦉%d. %s %s 🚀 %s%.2f%% \n\n💰 P/L: *$%.2f\n*💵 Collateral: $%s\n%s %sx\n🔗 Polygon\n🎯 Gains Network\n\nEntry Price: $%.2f\nCurrent Price: $%.2f\nTP: $%s SL: $%s Liq: $%.2f\nPotential Payout at TP: $%.2f\n\n",
		activeIdxIncr,
		tradeDirection,
		ticker,
		getPnlOperator,
		getPnlPrcnt,
		getPnlUSD,
		activeTrade.PositionSizeDai,
		tradeDirection,
		activeTrade.Leverage,
		openPrice,
		currentPrice,
		activeTrade.TP,
		activeTrade.SL,
		liqPrice,
		payOut,
	)

	var otherPositionStr string
	// Iterate over the array
	for i := 0; i < len(openTradesArraySorted); i++ {
		// Skip if we hit activeIdx
		if i == activeIdx {
			continue
		}
		openTradesStr, err := BuildSinglePositionString(rdbPrice, openTradesArraySorted[i], i)
		if err != nil {
			log.Println("Error getting generating position string in BuildPositionsAndOrderString")
		}
		otherPositionStr += openTradesStr + "\n"
	}

	return activePositionString + otherPositionStr
}
